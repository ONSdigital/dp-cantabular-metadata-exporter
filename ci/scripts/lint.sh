#!/bin/bash -eux

pushd dp-cantabular-metadata-exporter
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
  make lint
popd
