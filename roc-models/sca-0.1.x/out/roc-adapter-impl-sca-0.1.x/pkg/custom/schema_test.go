// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

package custom

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadJson(t *testing.T) {
	content, err := os.ReadFile("testdata/sample-config.json")
	assert.Nil(t, err)

	var payload sraConfig
	err = json.Unmarshal(content, &payload)
	assert.Nil(t, err)

	assert.Equal(t, payload.Version, float32(1.0))

	assert.Equal(t, payload.Models["person_vehicle_bike_model"][0].Name, "person-vehicle-bike-detection-crossroad-0078")
	assert.Equal(t, payload.Models["person_vehicle_bike_model"][0].Description, "https://github.com/openvinotoolkit/open_model_zoo/blob/master/models/intel/person-vehicle-bike-detection-crossroad-0078/README.md")
	assert.Equal(t, payload.Models["person_vehicle_bike_model"][0].Size, "3.964GFlops")

	assert.Equal(t, payload.Pipelines[0].Name, "person-vehicle-bike-detection")
	assert.Equal(t, payload.Pipelines[0].Info, "Application to detect person and vehicles")
	assert.Equal(t, payload.Pipelines[0].Enable, boolString("True"))
	assert.Equal(t, len(payload.Pipelines[0].Sources), 1)
	assert.Equal(t, payload.Pipelines[0].Sources[0].Name, "personvehicle-input")
	assert.Equal(t, payload.Pipelines[0].Sources[0].StreamCount, uint8(1))
	assert.Equal(t, len(payload.Pipelines[0].Nodes), 1)
	assert.Equal(t, payload.Pipelines[0].Nodes[0].Name, "person_vehicle_bike_model")
	assert.Equal(t, payload.Pipelines[0].Nodes[0].Model, "person-vehicle-bike-detection-crossroad-0078")
	assert.Equal(t, payload.Pipelines[0].Nodes[0].Precision, "FP16")
	assert.Equal(t, payload.Pipelines[0].Nodes[0].Device, "CPU")
}
