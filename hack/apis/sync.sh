#!/bin/bash -l
#
# Copyright 2021. The KubeVela Authors.
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
#

set -e

echo "echo VELA_API_DEPLOY"
echo $VELA_API_DEPLOY
echo "echo VELA_API_DEPLOY"

if [[ -n "$VELA_API_DEPLOY" ]]
then
  mkdir -p ~/.ssh
  echo "$VELA_API_DEPLOY" > ~/.ssh/id_rsa
  chmod 600 ~/.ssh/id_rsa
fi

echo "git clone"
cd ..
pwd
git config --global user.email "kubevela.bot@aliyun.com"
git config --global user.name "kubevela-bot"
git clone --single-branch --depth 1 git@github.com:yangsoon/kubevela-core-api.git kubevela-core-api

pwd
echo "clear kubevela-core-api api/"
rm -r kubevela-core-api/apis/*

echo "clear kubevela-core-api pkg/oam"
rm -r kubevela-core-api/pkg/oam/*

echo "update kubevela-core-api api/"
cp -R kubevela/apis/* kubevela-core-api/apis/

echo "update kubevela-core-api pkg/oam"
cp -R kubevela/pkg/oam/* kubevela-core-api/pkg/oam/

echo "change import path"
find ./kubevela-core-api -type f -name "*.go" -print0 | xargs -0 sed -i '' 's|github.com/oam-dev/kubevela/|github.com/oam-dev/kubevela-core-api/|g'

echo "test api"
cd kubevela-core-api
pwd
go build test/main.go

echo "push to kubevela-core-api"
if git diff --quiet
then
  echo "nothing need to push, finished!"
else
  git add .
  git commit -m "sync commit $COMMIT_ID from kubevela-$VERSION"
  git push origin main
fi