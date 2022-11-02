# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: LicenseRef-Intel

FROM golang:1.19-alpine AS build

RUN mkdir roc-api-server
WORKDIR /roc-api-server

COPY ./go.mod /roc-api-server
COPY ./go.sum /roc-api-server
COPY ./cmd /roc-api-server/cmd
COPY ./pkg /roc-api-server/pkg
COPY ./vendor /roc-api-server/vendor

RUN go build -mod=vendor -o ./roc-rest-api ./cmd/main.go


FROM alpine:3.16

COPY --from=build /roc-api-server/roc-rest-api /usr/local/bin/roc-rest-api
COPY ./assets /assets

ENTRYPOINT ["roc-rest-api"]
