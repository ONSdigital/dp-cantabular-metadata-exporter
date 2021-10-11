#!/bin/bash -eux

pushd dp-cantabular-metadata-exporter
  make build
  cp build/dp-cantabular-metadata-exporter Dockerfile.concourse ../build
popd
