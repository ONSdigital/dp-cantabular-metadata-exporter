#!/bin/bash -eux

go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/ONSdigital/dis-vulncheck@v0.6.1

pushd dp-cantabular-metadata-exporter
  make audit
popd