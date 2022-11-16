// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

// Package custom implements a synchronizer for converting sra gnmi to custom
// format
package custom

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/gnmi"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var (
	StaticConfigFileName = flag.String("static_config_file", "config.json", "Static configuration for pipelines and models")
	MockPushUrl          *string
)

var log = logging.GetLogger("custom")

func (ss *SCASyncStep) SynchronizeSources() {

	// be deterministic
	districtKeys := []string{}
	for k := range ss.City.District {
		districtKeys = append(districtKeys, k)
	}
	sort.Strings(districtKeys)

	for _, districtKey := range districtKeys {
		district := ss.City.District[districtKey]

		// be deterministic
		srcKeys := []string{}
		for k := range district.Source {
			srcKeys = append(srcKeys, k)
		}
		sort.Strings(srcKeys)

		for _, srcKey := range srcKeys {
			src := district.Source[srcKey]
			source := sourceDefinition{
				Name: *src.SourceId,
			}
			if src.Video != nil {
				source.Source = *src.Video.Path
				source.Type = src.Video.SourceType.String()
			}
			ss.Config.Sources = append(ss.Config.Sources, source)
		}
	}
}

func (ss *SCASyncStep) SynchronizeModels() {
	ss.Config.Models = ss.StaticConfig.Models
}

func (ss *SCASyncStep) SynchronizePipeline(srcPipeline *pipelineDefinition) error {
	// Copy pipeline attributes from source JSON
	pipeline := pipelineDefinition{
		Name: srcPipeline.Name,
		Info: srcPipeline.Info,
	}

	// Get the enable state for this pipeline from ROC
	enable, err := ss.LookupEnable(&srcPipeline.Name)
	if err != nil {
		return err
	}
	if enable {
		pipeline.Enable = "True"
	} else {
		pipeline.Enable = "False"
	}

	// Copy Nodes from the source JSON
	pipeline.Nodes = srcPipeline.Nodes

	for k, node := range pipeline.Nodes {
		model, precision, device, err := ss.LookupApplication(&pipeline.Name, &node.Name)
		if err != nil {
			return err
		}
		if model != nil {
			pipeline.Nodes[k].Model = *model
		}
		if precision != nil {
			pipeline.Nodes[k].Precision = *precision
		}
		if device != nil {
			pipeline.Nodes[k].Device = *device
		}
	}

	// Resolve all the district references for this pipeline
	districtRefs, err := ss.LookupDistricts(&srcPipeline.Name)
	if err != nil {
		return err
	}

	// For each district reference, get the list of sources from ROC and add them to the pipeline's source list
	for _, districtRef := range districtRefs {
		district, err := ss.LookupDistrict(&districtRef.DistrictID)
		if err != nil {
			return err
		}

		// be deterministic
		srcKeys := []string{}
		for k := range district.Source {
			srcKeys = append(srcKeys, k)
		}
		sort.Strings(srcKeys)

		for _, srcKey := range srcKeys {
			source := district.Source[srcKey]
			sr := sourceRef{Name: *source.SourceId,
				StreamCount: districtRef.StreamCount}
			pipeline.Sources = append(pipeline.Sources, sr)
		}
	}

	ss.Config.Pipelines = append(ss.Config.Pipelines, pipeline)

	return nil
}

func (ss *SCASyncStep) SynchronizeTarget() (int, error) {
	pushFailures := 0

	// Synchronize the list of sources
	ss.SynchronizeSources()

	// Synchronize the list of models
	ss.SynchronizeModels()

	// Now do the pipelines. The source JSON contains the list of pipelines and the set
	// of Nodes for each pipeline. So we'll iterate through that, and then fill in ROC data
	// for the sources and models that feed into the pipeline.
	for _, pipeline := range ss.StaticConfig.Pipelines {
		err := ss.SynchronizePipeline(&pipeline)
		if err != nil {
			return pushFailures, err
		}
	}

	// this could happen if the static json is missing pipelines for the retail store
	if len(ss.Config.Pipelines) == 0 {
		log.Infof("Skipping target %s because it has no pipelines", *ss.CityID)
		return 0, nil
	}

	// this could happen if the ROC config does not have any source setup for the store
	if len(ss.Config.Sources) == 0 {
		log.Infof("Skipping target %s because it has no sources", *ss.CityID)
		return 0, nil
	}

	if ss.Synchronizer.IsPartialUpdateEnabled() && ss.Synchronizer.CacheCheck(CacheModelConfig, *ss.CityID, ss.Config) {
		log.Infof("Store %s config has not changed", *ss.CityID)
		return 0, nil
	}

	data, err := json.MarshalIndent(ss.Config, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("store %s failed to Marshal config Json: %s", *ss.CityID, err)
	}

	if (ss.Endpoint == nil) || (*ss.Endpoint == "") {
		return 0, fmt.Errorf("store %s has no netconfig endpoint to push to", *ss.CityID)
	}

	url := *ss.Endpoint
	restPusher := synchronizer.NewRestPusher(url, data, synchronizer.WithContentType("application/json"))
	err = restPusher.PushUpdate()
	if err != nil {
		return 1, fmt.Errorf("store %s failed to Push config update: %s", *ss.CityID, err)
	}

	ss.Synchronizer.CacheUpdate(CacheModelConfig, *ss.CityID, ss.Config)

	log.Infof("post: %s", data)

	return pushFailures, nil
}

// SynchronizeDevice synchronizes a device. Two sets of error state are returned:
//  1. pushFailures -- a count of pushes that failed to the core. Synchronizer should retry again later.
//  2. error -- a fatal error that occurred during synchronization.
func SynchronizeDevice(ctx context.Context, sync synchronizer.SynchronizerInterface, allConfig *gnmi.ConfigForest) (int, error) {
	// Read staticConfig, aka the "source json"
	// This is the configuration that will be generated by springboard that is not overridden by the ROC. It's currently the
	// same format as our output JSON, so we can read it into the schema.

	staticConfigContent, err := os.ReadFile(*StaticConfigFileName)
	if err != nil {
		return 0, fmt.Errorf("failed to read static config files %s: %s", *StaticConfigFileName, err)
	}

	var staticConfig sraConfig
	err = json.Unmarshal(staticConfigContent, &staticConfig)
	if err != nil {
		return 0, fmt.Errorf("failed to parse static config files %s: %s", *StaticConfigFileName, err)
	}

	pushFailures := 0
	for targetID, targetConfig := range allConfig.Configs {
		device := targetConfig.(*SmartCity)

		ss := &SCASyncStep{
			CityID:       &targetID,
			City:         device,
			Config:       &sraConfig{},
			StaticConfig: &staticConfig,
			Synchronizer: sync,
		}

		if MockPushUrl != nil {
			ss.Endpoint = MockPushUrl
		} else {
			controllerInfo, err := lookupSRAControllerInfo(ctx, sync, targetID)
			if err != nil {
				log.Warnf("Failed to resolve controller info for Store %s", targetID)
				// leave the endpoint blank
			} else {
				var endpoint string
				if (strings.Contains(controllerInfo.ControlEndpoint.Address, ":")) || (controllerInfo.ControlEndpoint.Port == 0) {
					// The address already contains a ":", which implies that it has a port number inside of it, or the
					// port number is 0. Either way it's clear the caller doesn't want us inserting a port number into
					// what has been passed. 'Address' could be a partial URL, for example hostname/foo. We do this when using
					// nginx reverse proxy.
					endpoint = fmt.Sprintf("http://%s/", controllerInfo.ControlEndpoint.Address)
				} else {
					endpoint = fmt.Sprintf("http://%s:%d/", controllerInfo.ControlEndpoint.Address, controllerInfo.ControlEndpoint.Port)
				}
				log.Infof("Endpoint from topo is %s", endpoint)
				ss.Endpoint = &endpoint
			}
		}

		ss.Config.Version = ss.StaticConfig.Version

		thisPushFailures, err := ss.SynchronizeTarget()
		pushFailures += thisPushFailures
		if err != nil {
			// Adapter semantics are to log the error and continue synchronizing other targets. If pushFailures>1, then the
			// adapter will come around and try to hit it again.
			log.Warnf("Target %s failed to synchronize: %s", targetID, err)
		}

	}

	return pushFailures, nil
}
