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

	ss := &SRASyncStep{
		StoreID: synchronizer.AStr("foo"),
		Store:   device,
	}

	_, err := ss.LookupEnable(synchronizer.AStr("nonexistent"))
	assert.EqualError(t, err, "unknown pipeline nonexistent")

	_, err = ss.LookupEnable(nil)
	assert.EqualError(t, err, "nil pipeline")

	var b bool

	b, err = ss.LookupEnable(synchronizer.AStr("shopper-pose-mood-estimation"))
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	b, err = ss.LookupEnable(synchronizer.AStr("shopper-count-duration"))
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	b, err = ss.LookupEnable(synchronizer.AStr("shelf-object-count"))
	assert.Nil(t, err)
	assert.Equal(t, true, b)
}

func TestLookupArea(t *testing.T) {
	_, device := buildSampleConfig()

	ss := &SRASyncStep{
		StoreID: synchronizer.AStr("foo"),
		Store:   device,
	}

	_, err := ss.LookupArea(synchronizer.AStr("nonexistent"))
	assert.EqualError(t, err, "failed to find retailarea nonexistent")

	var ra *RetailArea

	ra, err = ss.LookupArea(synchronizer.AStr("traf"))
	assert.Nil(t, err)
	assert.Equal(t, "traf", *ra.AreaId)
}
