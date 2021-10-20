#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-cantabular-metadata-exporter
  make lint
popd
