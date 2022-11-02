// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

// Package custom implements a synchronizer for converting sra gnmi to custom
// format
package custom

import (
	"context"
	"encoding/json"
	"fmt"
    models "github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-models/models/sra-0.2.x/api"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/gnmi"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"
)

// SynchronizeDevice synchronizes a device. Two sets of error state are returned:
//   1) pushFailures -- a count of pushes that failed to the core. Synchronizer should retry again later.
//   2) error -- a fatal error that occurred during synchronization.
func SynchronizeDevice(ctx context.Context, sync synchronizer.SynchronizerInterface, allConfig *gnmi.ConfigForest) (int, error) {
	_ = sync

	pushFailures := 0
	for targetID, targetConfig := range allConfig.Configs {
		device := targetConfig.(*models.Device)
		if device != nil {
			fmt.Println("Target ID : ", targetID)
			//spew.Dump(device)
			deviceJSON, err := json.MarshalIndent(device, "", "  ")
			if err != nil {
				fmt.Errorf(err.Error())
			}
			fmt.Printf("%s\n", string(deviceJSON))
		}
	}

	// TODO: everything

	return pushFailures, nil
}
