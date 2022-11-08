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

	assert.Equal(t, payload.Models["face_detection"][0].Name, "face-detection-adas-0001")
	assert.Equal(t, payload.Models["face_detection"][0].Description, "https://github.com/openvinotoolkit/open_model_zoo/blob/master/models/intel/face-detection-adas-0001/README.md")
	assert.Equal(t, payload.Models["face_detection"][0].Size, "2.835GFlops")

	assert.Equal(t, payload.Pipelines[0].Name, "shopper-pose-mood-estimation")
	assert.Equal(t, payload.Pipelines[0].Info, "Application to recognize Shopper Emotions")
	assert.Equal(t, payload.Pipelines[0].Enable, boolString("True"))
	assert.Equal(t, len(payload.Pipelines[0].Sources), 2)
	assert.Equal(t, payload.Pipelines[0].Sources[0].Name, "headpose-input1")
	assert.Equal(t, payload.Pipelines[0].Sources[0].StreamCount, uint8(2))
	assert.Equal(t, len(payload.Pipelines[0].Nodes), 3)
	assert.Equal(t, payload.Pipelines[0].Nodes[0].Name, "face_detection")
	assert.Equal(t, payload.Pipelines[0].Nodes[0].Model, "face-detection-adas-0001")
	assert.Equal(t, payload.Pipelines[0].Nodes[0].Precision, "FP16")
	assert.Equal(t, payload.Pipelines[0].Nodes[0].Device, "CPU")
}
