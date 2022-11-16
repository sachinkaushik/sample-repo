// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

package custom

import (
	"context"

	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/gnmi"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"
	models "github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sca-0.1.x/api"
)

func emptySynchronizeDevice(ctx context.Context, sync synchronizer.SynchronizerInterface, config *gnmi.ConfigForest) (int, error) {
	_ = ctx
	_ = sync
	_ = config
	return 0, nil
}

func newMockSynchronizer() (synchronizer.SynchronizerInterface, error) {
	return synchronizer.NewSynchronizer(emptySynchronizeDevice)
}

func buildSampleSmartCity() *SmartCity {
	src1 := Source{
		Description: synchronizer.AStr("Source1Desc"),
		DisplayName: synchronizer.AStr("Source1DisplayName"),
		Image:       synchronizer.AStr("Source1Image"),
		SourceId:    synchronizer.AStr("source1"),
		Video:       &SourceVideo{Path: synchronizer.AStr("/path/to/source1"), SourceType: VideoSample},
	}

	src2 := Source{
		Description: synchronizer.AStr("Source2Desc"),
		DisplayName: synchronizer.AStr("Source2DisplayName"),
		Image:       synchronizer.AStr("Source2Image"),
		SourceId:    synchronizer.AStr("source2"),
		Video:       &SourceVideo{Path: synchronizer.AStr("/path/to/source2"), SourceType: VideoSample},
	}

	src3 := Source{
		Description: synchronizer.AStr("Source3Desc"),
		DisplayName: synchronizer.AStr("Source3DisplayName"),
		Image:       synchronizer.AStr("Source3Image"),
		SourceId:    synchronizer.AStr("source3"),
		Video:       &SourceVideo{Path: synchronizer.AStr("/path/to/source3"), SourceType: VideoSample},
	}

	ra1 := District{
		DistrictId:  synchronizer.AStr("headpose"),
		Description: synchronizer.AStr("headposedesc"),
		DisplayName: synchronizer.AStr("headposedn"),
		Source: map[string]*models.IntelScaSource_District_Source{
			*src1.SourceId: &src1,
		},
	}

	ra2 := District{
		DistrictId:  synchronizer.AStr("objdetect"),
		Description: synchronizer.AStr("objdetectdesc"),
		DisplayName: synchronizer.AStr("objdetectdn"),
		Source: map[string]*models.IntelScaSource_District_Source{
			*src2.SourceId: &src2,
		},
	}

	ra3 := District{
		DistrictId:  synchronizer.AStr("traf"),
		Description: synchronizer.AStr("trafdesc"),
		DisplayName: synchronizer.AStr("trafdn"),
		Source: map[string]*models.IntelScaSource_District_Source{
			*src3.SourceId: &src3,
		},
	}

	coll_ra := models.IntelSca_CollisionDetection_District{
		DistrictRef: synchronizer.AStr(*ra1.DistrictId),
		Enabled:     synchronizer.ABool(true),
		StreamCount: synchronizer.AUint8(1),
	}

	coll := CollisionDetection{
		District: map[string]*models.IntelSca_CollisionDetection_District{"shop_cams": &coll_ra},
		Enable:   synchronizer.ABool(true),
	}

	class_ra := models.IntelSca_TrafficClassification_District{
		DistrictRef: synchronizer.AStr(*ra2.DistrictId),
		Enabled:     synchronizer.ABool(true),
		StreamCount: synchronizer.AUint8(1),
	}

	class := TrafficClassification{
		District: map[string]*models.IntelSca_TrafficClassification_District{"shelf_cams": &class_ra},
		Enable:   synchronizer.ABool(true),
	}

	traf_ra := models.IntelSca_TrafficMonitoring_District{
		DistrictRef: synchronizer.AStr(*ra2.DistrictId),
		Enabled:     synchronizer.ABool(true),
		StreamCount: synchronizer.AUint8(1),
	}

	traf := TrafficMonitoring{
		District: map[string]*models.IntelSca_TrafficMonitoring_District{"traf_cams": &traf_ra},
		Enable:   synchronizer.ABool(true),
	}

	device := &SmartCity{
		District: map[string]*District{
			*ra1.DistrictId: &ra1,
			*ra2.DistrictId: &ra2,
			*ra3.DistrictId: &ra3,
		},
		CollisionDetection:    &coll,
		TrafficClassification: &class,
		TrafficMonitoring:     &traf,
	}

	return device
}

func buildSampleConfig() (*gnmi.ConfigForest, *SmartCity) {
	device := buildSampleSmartCity()
	config := gnmi.NewConfigForest()
	config.Configs["sample-ent"] = device
	return config, device
}
