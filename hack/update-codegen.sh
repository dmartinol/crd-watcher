#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

PROJECT_MODULE="github.com/dmartinol/crd-watcher"
CUSTOM_RESOURCE_NAME="requeststate"
CUSTOM_RESOURCE_VERSION="v1"

# go mod vendor
chmod +x "$SCRIPT_ROOT"/vendor/k8s.io/code-generator/generate-groups.sh

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
"${CODEGEN_PKG}/generate-groups.sh" "all" \
  "$PROJECT_MODULE/pkg/client" \
  "$PROJECT_MODULE/pkg/apis" \
  "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt
  # --output-base "$(dirname "${BASH_SOURCE[0]}")/../pkg"

  mv ${SCRIPT_ROOT}/$PROJECT_MODULE/pkg/apis/requeststate/v1/* "${SCRIPT_ROOT}"/pkg/apis/requeststate/v1
  rm -rf "${SCRIPT_ROOT}"/pkg/client
  mv $PROJECT_MODULE/pkg/client "${SCRIPT_ROOT}"/pkg
  rm -rf "${SCRIPT_ROOT}/$PROJECT_MODULE"
