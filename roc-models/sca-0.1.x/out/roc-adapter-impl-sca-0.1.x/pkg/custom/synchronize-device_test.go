// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

package custom

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/gnmi"
	"github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var TestDataSampleConfigFileName = "testdata/sample-config.json"

func TestSynchronizeDeviceEmpty(t *testing.T) {
	sync, err := newMockSynchronizer()
	assert.Nil(t, err)

	StaticConfigFileName = &TestDataSampleConfigFileName
	allConfig := gnmi.ConfigForest{}
	pushErrors, err := SynchronizeDevice(context.Background(), sync, &allConfig)

	assert.Nil(t, err)
	assert.Equal(t, 0, pushErrors)
}

func TestSynchronizeDeviceNoStaticConfig(t *testing.T) {
	postCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		// TODO: insert assert here
		_ = body
		postCount += 1
	}))
	defer ts.Close()

	// this will prevent synchronizer from calling topo and hanging
	MockPushUrl = &ts.URL

	sync, err := newMockSynchronizer()
	assert.Nil(t, err)

	StaticConfigFileName = synchronizer.AStr("nonexistent-file")
	allConfig, _ := buildSampleConfig()

	_, err = SynchronizeDevice(context.Background(), sync, allConfig)

	assert.EqualError(t, err, "failed to read static config files nonexistent-file: open nonexistent-file: no such file or directory")
}

func TestSynchronizeDeviceNoSources(t *testing.T) {
	postCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		// TODO: insert assert here
		_ = body
		postCount += 1
	}))
	defer ts.Close()

	// this will prevent synchronizer from calling topo and hanging
	MockPushUrl = &ts.URL

	sync, err := newMockSynchronizer()
	assert.Nil(t, err)

	StaticConfigFileName = &TestDataSampleConfigFileName
	allConfig, device := buildSampleConfig()

	device.District = map[string]*District{}
	pushErrors, err := SynchronizeDevice(context.Background(), sync, allConfig)

	assert.Nil(t, err)
	assert.Equal(t, 0, pushErrors)
	assert.Equal(t, 0, postCount)
}

func TestSynchronizeDeviceNoApps(t *testing.T) {
	pushes := map[string]string{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		pushes[r.URL.String()] = string(body)
	}))
	defer ts.Close()

	// this will prevent synchronizer from calling topo and hanging
	MockPushUrl = &ts.URL

	sync, err := newMockSynchronizer()
	assert.Nil(t, err)

	StaticConfigFileName = &TestDataSampleConfigFileName
	allConfig, device := buildSampleConfig()

	device.CollisionDetection = nil
	device.TrafficClassification = nil
	device.TrafficMonitoring = nil
	pushErrors, err := SynchronizeDevice(context.Background(), sync, allConfig)

	assert.Nil(t, err)
	assert.Equal(t, 0, pushErrors)
	assert.Equal(t, 1, len(pushes))

	jsonData, err := os.ReadFile("./testdata/out-syncdevice-noapps.json")
	assert.NoError(t, err)

	json, okay := pushes["/"]
	assert.True(t, okay)
	if okay {
		require.JSONEq(t, string(jsonData), json)
	}
}

func TestSynchronizeDevice(t *testing.T) {
	pushes := map[string]string{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		pushes[r.URL.String()] = string(body)
	}))
	defer ts.Close()

	// this will prevent synchronizer from calling topo and hanging
	MockPushUrl = &ts.URL

	sync, err := newMockSynchronizer()
	assert.Nil(t, err)

	StaticConfigFileName = &TestDataSampleConfigFileName
	allConfig, _ := buildSampleConfig()

	pushErrors, err := SynchronizeDevice(context.Background(), sync, allConfig)

	assert.Nil(t, err)
	assert.Equal(t, 0, pushErrors)
	assert.Equal(t, 1, len(pushes))

	jsonData, err := os.ReadFile("./testdata/out-syncdevice.json")
	assert.NoError(t, err)

	json, okay := pushes["/"]
	assert.True(t, okay)
	if okay {
		require.JSONEq(t, string(jsonData), json)
	}
}
