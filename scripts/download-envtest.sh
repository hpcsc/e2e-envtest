#!/bin/bash

set -euo pipefail

ENVTEST_VERSION=1.21.2
ENVTEST_DIR=./bin/envtest

if [ -d ${ENVTEST_DIR} ]; then
  echo "=== envtest exists at ${ENVTEST_DIR}, skip downloading"
else
  mkdir -p ${ENVTEST_DIR}
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  curl -sSL "https://go.kubebuilder.io/test-tools/${ENVTEST_VERSION}/${OS}/amd64" | tar -C ${ENVTEST_DIR} --strip-components=2 -xvzf -
  echo "=== downloaded envtest ${ENVTEST_VERSION} to ${ENVTEST_DIR}"
fi
