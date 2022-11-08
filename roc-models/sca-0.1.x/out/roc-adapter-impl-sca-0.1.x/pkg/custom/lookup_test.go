// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

package custom

import (
	"testing"

	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"
	"github.com/stretchr/testify/assert"
)

func TestLookupEnable(t *testing.T) {
	_, device := buildSampleConfig()

	ss := &SCASyncStep{
		CityID: synchronizer.AStr("foo"),
		City:   device,
	}

	_, err := ss.LookupEnable(synchronizer.AStr("nonexistent"))
	assert.EqualError(t, err, "unknown pipeline nonexistent")

	_, err = ss.LookupEnable(nil)
	assert.EqualError(t, err, "nil pipeline")

	var b bool

	b, err = ss.LookupEnable(synchronizer.AStr(PipeLineCollisionDetection))
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	b, err = ss.LookupEnable(synchronizer.AStr(PipelineTrafficClassification))
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	b, err = ss.LookupEnable(synchronizer.AStr(PipelineTrafficMonitoring))
	assert.Nil(t, err)
	assert.Equal(t, true, b)
}

func TestLookupArea(t *testing.T) {
	_, device := buildSampleConfig()

	ss := &SCASyncStep{
		CityID: synchronizer.AStr("foo"),
		City:   device,
	}

	_, err := ss.LookupDistrict(synchronizer.AStr("nonexistent"))
	assert.EqualError(t, err, "failed to find district nonexistent")

	var ra *District

	ra, err = ss.LookupDistrict(synchronizer.AStr("traf"))
	assert.Nil(t, err)
	assert.Equal(t, "traf", *ra.DistrictId)
}
