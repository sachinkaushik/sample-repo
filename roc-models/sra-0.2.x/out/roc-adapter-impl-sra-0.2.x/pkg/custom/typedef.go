// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

// Package custom implements a synchronizer for converting sra gnmi to custom
// format
package custom

import (
	models "github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sra-0.2.x/api"
)

const (
	VideoUnset  = models.IntelSraSource_RetailArea_Source_Video_SourceType_UNSET
	VideoFile   = models.IntelSraSource_RetailArea_Source_Video_SourceType_file
	VideoStream = models.IntelSraSource_RetailArea_Source_Video_SourceType_stream
	VideoDevice = models.IntelSraSource_RetailArea_Source_Video_SourceType_device
	VideoSample = models.IntelSraSource_RetailArea_Source_Video_SourceType_sample

	PipelineShopperMonitoring      = "shopper-pose-mood-estimation"
	PipelineStoreTrafficMonitoring = "shopper-count-duration"
	PipelineShelfMonitoring        = "shelf-object-count"

	AppObjectDetection    = "object_detection"
	AppFaceDetection      = "face_detection"
	AppPoseEstimation     = "pose_estimation"
	AppEmotionRecognition = "emotion_recognition"
	AppPersonDetection    = "person_detection"
)

type RetailStore = models.Device                                       //nolint
type RetailArea = models.IntelSraSource_RetailArea                     //nolint
type ShelfMonitoring = models.IntelSra_ShelfMonitoring                 //nolint
type ShopperMonitoring = models.IntelSra_ShopperMonitoring             //nolint
type StoreTrafficMonitoring = models.IntelSra_StoreTrafficMonitoring   //nolint
type RetailAreaLocation = models.IntelSraSource_RetailArea_Location    //nolint
type Source = models.IntelSraSource_RetailArea_Source                  //nolint
type SourceLocation = models.IntelSraSource_RetailArea_Source_Location //nolint
type SourceState = models.IntelSraSource_RetailArea_Source_State       //nolint
type SourceVideo = models.IntelSraSource_RetailArea_Source_Video       //nolint
