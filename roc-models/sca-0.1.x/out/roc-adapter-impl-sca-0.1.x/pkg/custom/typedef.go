// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

// Package custom implements a synchronizer for converting sra gnmi to custom
// format
package custom

import (
	models "github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sca-0.1.x/api"
)

const (
	VideoUnset  = models.IntelScaSource_District_Source_Video_SourceType_UNSET
	VideoFile   = models.IntelScaSource_District_Source_Video_SourceType_file
	VideoStream = models.IntelScaSource_District_Source_Video_SourceType_stream
	VideoDevice = models.IntelScaSource_District_Source_Video_SourceType_device
	VideoSample = models.IntelScaSource_District_Source_Video_SourceType_sample

	PipeLineCollisionDetection    = "vehicle-collision-detection"
	PipelineTrafficClassification = "vehicle-classification"
	PipelineTrafficMonitoring     = "person-vehicle-bike-detection"
)

type SmartCity = models.Device                                       //nolint
type District = models.IntelScaSource_District                       //nolint
type CollisionDetection = models.IntelSca_CollisionDetection         //nolint
type TrafficClassification = models.IntelSca_TrafficClassification   //nolint
type TrafficMonitoring = models.IntelSca_TrafficMonitoring           //nolint
type DistrictLocation = models.IntelScaSource_District_Location      //nolint
type Source = models.IntelScaSource_District_Source                  //nolint
type SourceLocation = models.IntelScaSource_District_Source_Location //nolint
type SourceState = models.IntelScaSource_District_Source_State       //nolint
type SourceVideo = models.IntelScaSource_District_Source_Video       //nolint
