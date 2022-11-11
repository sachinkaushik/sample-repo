// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

package custom

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	onoserrors "github.com/onosproject/onos-lib-go/pkg/errors"
)

// LookupEnable given a pipeline name, returns the enable status
func (ss *SRASyncStep) LookupEnable(pipeline *string) (bool, error) {
	if pipeline == nil {
		return false, fmt.Errorf("nil pipeline")
	}
	if *pipeline == PipelineShopperMonitoring {
		return (ss.Store != nil) &&
			(ss.Store.ShopperMonitoring != nil) &&
			(ss.Store.ShopperMonitoring.Enable != nil) && *ss.Store.ShopperMonitoring.Enable, nil
	} else if *pipeline == PipelineStoreTrafficMonitoring {
		return (ss.Store != nil) &&
			(ss.Store.StoreTrafficMonitoring != nil) &&
			(ss.Store.StoreTrafficMonitoring.Enable != nil) && *ss.Store.StoreTrafficMonitoring.Enable, nil
	} else if *pipeline == PipelineShelfMonitoring {
		return (ss.Store != nil) &&
			(ss.Store.ShelfMonitoring != nil) &&
			(ss.Store.ShelfMonitoring.Enable != nil) && *ss.Store.ShelfMonitoring.Enable, nil
	} else {
		return false, fmt.Errorf("unknown pipeline %s", *pipeline)
	}
}

// LookupApplication given a pipeline name and a node name, return details about the application
func (ss *SRASyncStep) LookupApplication(pipeline *string, node *string) (*string, *string, *string, error) {
	if pipeline == nil {
		return nil, nil, nil, fmt.Errorf("nil pipeline")
	}
	if node == nil {
		return nil, nil, nil, fmt.Errorf("nil node")
	}

	var model *string = nil
	var precision *string = nil
	var device *string = nil

	switch *pipeline {
	case PipelineShopperMonitoring:
		switch *node {
		case AppFaceDetection:
			if (ss.Store != nil) && (ss.Store.ShopperMonitoring != nil) && (ss.Store.ShopperMonitoring.FaceDetectionApplication != nil) {
				model = ss.Store.ShopperMonitoring.FaceDetectionApplication.Model
				precision = synchronizer.AStr(ss.Store.ShopperMonitoring.FaceDetectionApplication.Precision.String())
				device = synchronizer.AStr(ss.Store.ShopperMonitoring.FaceDetectionApplication.Device.String())
			}
		case AppPoseEstimation:
			if (ss.Store != nil) && (ss.Store.ShopperMonitoring != nil) && (ss.Store.ShopperMonitoring.HeadPoseDetectionApplication != nil) {
				model = ss.Store.ShopperMonitoring.HeadPoseDetectionApplication.Model
				precision = synchronizer.AStr(ss.Store.ShopperMonitoring.HeadPoseDetectionApplication.Precision.String())
				device = synchronizer.AStr(ss.Store.ShopperMonitoring.HeadPoseDetectionApplication.Device.String())
			}
		case AppEmotionRecognition:
			if (ss.Store != nil) && (ss.Store.ShopperMonitoring != nil) && (ss.Store.ShopperMonitoring.EmotionRecognitionApplication != nil) {
				model = ss.Store.ShopperMonitoring.EmotionRecognitionApplication.Model
				precision = synchronizer.AStr(ss.Store.ShopperMonitoring.EmotionRecognitionApplication.Precision.String())
				device = synchronizer.AStr(ss.Store.ShopperMonitoring.EmotionRecognitionApplication.Device.String())
			}
		default:
			return nil, nil, nil, fmt.Errorf("unknown node %s", *node)
		}
	case PipelineStoreTrafficMonitoring:
		switch *node {
		case AppPersonDetection:
			if (ss.Store != nil) && (ss.Store.StoreTrafficMonitoring != nil) && (ss.Store.StoreTrafficMonitoring.PersonDetectionApplication != nil) {
				model = ss.Store.StoreTrafficMonitoring.PersonDetectionApplication.Model
				precision = synchronizer.AStr(ss.Store.StoreTrafficMonitoring.PersonDetectionApplication.Precision.String())
				device = synchronizer.AStr(ss.Store.StoreTrafficMonitoring.PersonDetectionApplication.Device.String())
			}
		default:
			return nil, nil, nil, fmt.Errorf("unknown node %s", *node)
		}
	case PipelineShelfMonitoring:
		switch *node {
		case AppObjectDetection:
			if (ss.Store != nil) && (ss.Store.ShelfMonitoring != nil) && (ss.Store.ShelfMonitoring.ObjectDetectionApplication != nil) {
				model = ss.Store.ShelfMonitoring.ObjectDetectionApplication.Model
				precision = synchronizer.AStr(ss.Store.ShelfMonitoring.ObjectDetectionApplication.Precision.String())
				device = synchronizer.AStr(ss.Store.ShelfMonitoring.ObjectDetectionApplication.Device.String())
			}
		default:
			return nil, nil, nil, fmt.Errorf("unknown node %s", *node)
		}
	default:
		return nil, nil, nil, fmt.Errorf("unknown pipeline %s", *pipeline)
	}

	return model, precision, device, nil
}

// LookupArea given an area name returns that RetailArea
func (ss *SRASyncStep) LookupArea(ref *string) (*RetailArea, error) {
	if ref == nil {
		return nil, errors.New("retail area ref is nil")
	}

	ra, okay := ss.Store.RetailArea[*ref]
	if !okay {
		return nil, fmt.Errorf("failed to find retailarea %s", *ref)
	}

	return ra, nil
}

// LookupAreas given a pipeline name returns the set of area references for that pipeline
func (ss *SRASyncStep) LookupAreas(pipeline *string) ([]*AreaRef, error) {
	arearefs := []*AreaRef{}
	if pipeline == nil {
		return arearefs, fmt.Errorf("nil pipeline")
	}
	if *pipeline == PipelineShopperMonitoring {
		if ss.Store.ShopperMonitoring == nil {
			return arearefs, nil
		}
		// be deterministic
		areaKeys := []string{}
		for k := range ss.Store.ShopperMonitoring.RetailArea {
			areaKeys = append(areaKeys, k)
		}
		sort.Strings(areaKeys)

		for _, areaKey := range areaKeys {
			area := ss.Store.ShopperMonitoring.RetailArea[areaKey]
			if (area.Enabled != nil) && (*area.Enabled) {
				if area.AreaRef != nil && area.StreamCount != nil {
					arearefs = append(arearefs,
						&AreaRef{AreaID: *area.AreaRef, StreamCount: *area.StreamCount})
				}
			}
		}
		return arearefs, nil

	} else if *pipeline == PipelineStoreTrafficMonitoring {
		if ss.Store.StoreTrafficMonitoring == nil {
			return arearefs, nil
		}
		// be deterministic
		areaKeys := []string{}
		for k := range ss.Store.StoreTrafficMonitoring.RetailArea {
			areaKeys = append(areaKeys, k)
		}
		sort.Strings(areaKeys)

		for _, areaKey := range areaKeys {
			area := ss.Store.StoreTrafficMonitoring.RetailArea[areaKey]
			if (area.Enabled != nil) && (*area.Enabled) {
				if area.AreaRef != nil && area.StreamCount != nil {
					arearefs = append(arearefs,
						&AreaRef{AreaID: *area.AreaRef, StreamCount: *area.StreamCount})
				}
			}
		}
		return arearefs, nil

	} else if *pipeline == PipelineShelfMonitoring {
		if ss.Store.ShelfMonitoring == nil {
			return arearefs, nil
		}
		// be deterministic
		areaKeys := []string{}
		for k := range ss.Store.ShelfMonitoring.RetailArea {
			areaKeys = append(areaKeys, k)
		}
		sort.Strings(areaKeys)

		for _, areaKey := range areaKeys {
			area := ss.Store.ShelfMonitoring.RetailArea[areaKey]
			if (area.Enabled != nil) && (*area.Enabled) {
				if area.AreaRef != nil && area.StreamCount != nil {
					arearefs = append(arearefs,
						&AreaRef{AreaID: *area.AreaRef, StreamCount: *area.StreamCount})
				}
			}
		}
		return arearefs, nil
	}

	return arearefs, fmt.Errorf("unknown pipeline %s", *pipeline)
}

func lookupSRAControllerInfo(ctx context.Context, sync synchronizer.SynchronizerInterface, storeID string) (*topoapi.ControllerInfo, error) {
	topoClient, err := sync.GetTopoClient(ctx)
	if err != nil {
		return nil, onoserrors.FromGRPC(err)
	}

	getResponse, err := topoClient.Get(ctx, &topoapi.GetRequest{
		ID: topoapi.ID(storeID),
	})
	if err != nil {
		return nil, onoserrors.FromGRPC(err)
	}
	log.Debug("topo response object: %v", getResponse.Object)

	fabricObject := getResponse.Object
	controllerInfo := &topoapi.ControllerInfo{}
	err = fabricObject.GetAspect(controllerInfo)
	if err != nil {
		return nil, onoserrors.FromGRPC(err)
	}
	log.Debug("controller address %v port %v", controllerInfo.ControlEndpoint.Address, controllerInfo.ControlEndpoint.Port)

	return controllerInfo, nil
}
