// SPDX-FileCopyrightText: (C) 2022 Intel Corporation
// SPDX-License-Identifier: LicenseRef-Intel

// Package custom implements a synchronizer for converting sra gnmi to custom
// format
package custom

import "github.com/intel-innersource/frameworks.edge.one-intel-edge.maestro-app.roc.config-adapter/pkg/synchronizer"

// SRASyncStep holds the state for the custom sync step. We build it up as we go along.
type SRASyncStep struct {
	StoreID      *string
	Endpoint     *string
	Store        *RetailStore
	Config       *sraConfig
	StaticConfig *sraConfig
	Synchronizer synchronizer.SynchronizerInterface
}

// AreaRef is a reference to an area, obtained from ROC modeling.
type AreaRef struct {
	AreaID      string
	StreamCount uint8
}

const (
	// CacheModelConfig is the modelName to use when caching configs to SRA
	CacheModelConfig = "config"
)
