# SPDX-FileCopyrightText: (C) 2022 Intel Corporation
# SPDX-License-Identifier: LicenseRef-Intel

FROM golang:1.19-alpine as build

ENV ADAPTER_ROOT=$GOPATH/src/roc-adapter
ENV CGO_ENABLED=0

RUN mkdir -p $ADAPTER_ROOT/

COPY . $ADAPTER_ROOT/

RUN cat $ADAPTER_ROOT/go.mod
RUN apk add git

RUN cd $ADAPTER_ROOT && GO111MODULE=on go build -o /go/bin/roc-adapter \
        -ldflags \
        "-X roc-adapter/internal/pkg/version.Version=$org_label_schema_version \
         -X roc-adapter/internal/pkg/version.GitCommit=$org_label_schema_vcs_ref  \
         -X roc-adapter/internal/pkg/version.GitDirty=$org_opencord_vcs_dirty \
         -X roc-adapter/internal/pkg/version.GoVersion=$(go version 2>&1 | sed -E  's/.*go([0-9]+\.[0-9]+\.[0-9]+).*/\1/g') \
         -X roc-adapter/internal/pkg/version.Os=$(go env GOHOSTOS) \
         -X roc-adapter/internal/pkg/version.Arch=$(go env GOHOSTARCH) \
         -X roc-adapter/internal/pkg/version.BuildTime=$org_label_schema_build_date" \
         ./cmd/roc-adapter

FROM alpine:3.14
RUN apk add bash openssl curl libc6-compat

ENV HOME=/home/roc-adapter

RUN mkdir $HOME
WORKDIR $HOME

COPY --from=build /go/bin/roc-adapter /usr/local/bin/roc-adapter
COPY config/config.json /home/roc-adapter/
