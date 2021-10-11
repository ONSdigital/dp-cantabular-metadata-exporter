#!/bin/bash -eux

pushd dp-cantabular-metadata-exporter
  make test-component
popd
