// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/openconfig/ygot/ygot"

	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/gnmi"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/northbound"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/options"
	synchronizer "github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-adapters.sra/pkg/custom"
	models "github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sra-0.2.x/api"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("roc-adapter")

func getModels() *gnmi.Model {
	model := gnmi.NewModel(models.ModelData(),
		reflect.TypeOf((*models.Device)(nil)),
		models.SchemaTree["Device"],
		models.Unmarshal,
		//models.ΛEnum  // NOTE: There is no Enum in the SRA models? So use a blank map.
		map[string]map[int64]ygot.EnumDefinition{},
	)

	return model
}

func main() {
	// Load options and display banner
	options.LoadOptions("roc-adapter")

	// The synchronizer will convey its list of models.
	model := getModels()

	// If the ShowModelList option was used, then print out a
	// list of supported models and then exit.
	if *options.ShowModelList {
		fmt.Fprintf(os.Stdout, "Supported models:\n")
		for _, m := range model.SupportedModels() {
			fmt.Fprintf(os.Stdout, "  %s\n", m)
		}
		return
	}

	// Initialize the synchronizer's service-specific code.
	log.Infof("Initializing synchronizer")
	sync, err := synchronizer.NewSynchronizer(
		custom.SynchronizeDevice,
		synchronizer.WithPostEnable(!*options.PostDisable),
		synchronizer.WithPartialUpdateEnable(!*options.PartialUpdateDisable),
		synchronizer.WithPostTimeout(*options.PostTimeout),
		synchronizer.WithCertPaths(*options.CAPath, *options.KeyPath, *options.CertPath),
		synchronizer.WithTopoEndpoint(*options.TopoEndpoint),
	)

	if err != nil {
		panic(err)
	}

	sync.Start()

	// Setup the northbound APIs
	nb := northbound.NewNorthboundAPI(
		model,
		sync,
		northbound.WithBindAddr(*options.BindAddr),
		northbound.WithDiagsPort(*options.DiagsPort),
		northbound.WithOnosConfigAddr(*options.OnosConfigAddr),
		northbound.WithOnosConfigTarget(*options.OnosConfigTarget),
		northbound.WithMetricAddr(*options.MetricAddr),
	)
	nb.Start()
}
