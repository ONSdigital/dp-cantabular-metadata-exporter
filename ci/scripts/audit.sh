#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-cantabular-metadata-exporter
  make audit
popd