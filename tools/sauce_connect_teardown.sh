#!/bin/bash

# Copyright 2019 The Vitess Authors.
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

set -e -o pipefail

echo "Shutting down Sauce Connect tunnel"

killall sc

while [[ -n `ps -ef | grep "sauce-connect" | grep -v "grep"` ]]; do
  printf "."
  sleep .5
done

echo ""
echo "Sauce Connect tunnel has been shut down"
