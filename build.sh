#!/bin/bash

PACKAGE_NAME="prometheus_exporter"
MAIN_PACKAGE="cmd/${PACKAGE_NAME}/main.go"

BINARY_NAME="grl-exporter"

GITHASH=$(git rev-parse --short HEAD)
DATE=$(date -u)

CGO_ENABLED=0 GO111MODULE=auto go build -o ${BINARY_NAME} ${MAIN_PACKAGE}

zip ${BINARY_NAME}-${GOOS}-${GOARCH}.zip ${BINARY_NAME} > /dev/null
