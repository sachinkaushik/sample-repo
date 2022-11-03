#!/usr/bin/env bash
# SPDX-FileCopyrightText: (C) 2022 Intel Corporation
# SPDX-License-Identifier: LicenseRef-Intel

docker build --platform linux/x86_64 -t default-route-openshift-image-registry.apps.nex.one-edge.intel.com/springboard-dev-common/roc-sra-0.2.x:v0.0.1 roc-models/sra-0.2.x&

docker build --platform linux/x86_64 -f roc-models/sra-0.2.x/out/roc-api-server-sra-0.2.x/build/Dockerfile -t default-route-openshift-image-registry.apps.nex.one-edge.intel.com/springboard-dev-common/roc-api-server-sra-0.2.x:v0.0.1 roc-models/sra-0.2.x/out/roc-api-server-sra-0.2.x&

docker build --platform linux/x86_64 -f roc-models/sra-0.2.x/out/roc-adapter-stubs-sra-0.2.x/Dockerfile -t default-route-openshift-image-registry.apps.nex.one-edge.intel.com/springboard-dev-common/roc-adapter-stubs-sra-0.2.x:v0.0.1 roc-models/sra-0.2.x/out/roc-adapter-stubs-sra-0.2.x&
