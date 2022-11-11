// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

package custom

import (
	"context"

	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/gnmi"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"
	models "github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sra-0.2.x/api"
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

func buildSampleRetailStore() *RetailStore {
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

	ra1 := RetailArea{
		AreaId:      synchronizer.AStr("headpose"),
		Description: synchronizer.AStr("headposedesc"),
		DisplayName: synchronizer.AStr("headposedn"),
		Source: map[string]*models.IntelSraSource_RetailArea_Source{
			*src1.SourceId: &src1,
		},
	}

	ra2 := RetailArea{
		AreaId:      synchronizer.AStr("objdetect"),
		Description: synchronizer.AStr("objdetectdesc"),
		DisplayName: synchronizer.AStr("objdetectdn"),
		Source: map[string]*models.IntelSraSource_RetailArea_Source{
			*src2.SourceId: &src2,
		},
	}

	ra3 := RetailArea{
		AreaId:      synchronizer.AStr("traf"),
		Description: synchronizer.AStr("trafdesc"),
		DisplayName: synchronizer.AStr("trafdn"),
		Source: map[string]*models.IntelSraSource_RetailArea_Source{
			*src3.SourceId: &src3,
		},
	}

	shop_ra := models.IntelSra_ShopperMonitoring_RetailArea{
		AreaRef:     synchronizer.AStr(*ra1.AreaId),
		Enabled:     synchronizer.ABool(true),
		StreamCount: synchronizer.AUint8(1),
	}

	shop := ShopperMonitoring{
		RetailArea: map[string]*models.IntelSra_ShopperMonitoring_RetailArea{"shop_cams": &shop_ra},
		Enable:     synchronizer.ABool(true),
	}

	shelf_ra := models.IntelSra_ShelfMonitoring_RetailArea{
		AreaRef:     synchronizer.AStr(*ra2.AreaId),
		Enabled:     synchronizer.ABool(true),
		StreamCount: synchronizer.AUint8(1),
	}

	shelf := ShelfMonitoring{
		RetailArea: map[string]*models.IntelSra_ShelfMonitoring_RetailArea{"shelf_cams": &shelf_ra},
		Enable:     synchronizer.ABool(true),
	}

	traf_ra := models.IntelSra_StoreTrafficMonitoring_RetailArea{
		AreaRef:     synchronizer.AStr(*ra2.AreaId),
		Enabled:     synchronizer.ABool(true),
		StreamCount: synchronizer.AUint8(1),
	}

	traf := StoreTrafficMonitoring{
		RetailArea: map[string]*models.IntelSra_StoreTrafficMonitoring_RetailArea{"traf_cams": &traf_ra},
		Enable:     synchronizer.ABool(true),
	}

	device := &RetailStore{
		RetailArea: map[string]*RetailArea{
			*ra1.AreaId: &ra1,
			*ra2.AreaId: &ra2,
			*ra3.AreaId: &ra3,
		},
		ShelfMonitoring:        &shelf,
		ShopperMonitoring:      &shop,
		StoreTrafficMonitoring: &traf,
	}

	return device
}

func buildSampleConfig() (*gnmi.ConfigForest, *RetailStore) {
	device := buildSampleRetailStore()
	config := gnmi.NewConfigForest()
	config.Configs["sample-ent"] = device
	return config, device
}
