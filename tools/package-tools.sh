#!/bin/bash
set -e -x -o pipefail

BIN_DIR=${BIN_DIR:-./bin}

cd $(dirname $0)/..

if [ ! -e orion101-tools ]; then
    git clone --depth=1 https://github.com/orion101-ai/tools orion101-tools
fi

./orion101-tools/scripts/build.sh

for pj in $(find orion101-tools -name package.json | grep -v node_modules); do
    if [ $(basename $(dirname $pj)) == common ]; then
        continue
    fi
    (
        cd $(dirname $pj)
        echo Building $PWD
        pnpm i
    )
done

cd orion101-tools
if [ ! -e workspace-provider ]; then
    git clone --depth=1 https://github.com/gptscript-ai/workspace-provider
fi

cd workspace-provider
go build -ldflags="-s -w" -o bin/gptscript-go-tool .

cd ..

if [ ! -e datasets ]; then
    git clone --depth=1 https://github.com/gptscript-ai/datasets
fi

cd datasets
go build -ldflags="-s -w" -o bin/gptscript-go-tool .

cd ..

if [ ! -e aws-encryption-provider ]; then
    git clone --depth=1 https://github.com/kubernetes-sigs/aws-encryption-provider
fi

cd aws-encryption-provider
go build -o ${BIN_DIR}/aws-encryption-provider cmd/server/main.go

cd ../..

if ! command -v uv; then
    pip install uv
fi

if [ ! -e orion101-tools/venv ]; then
    uv venv orion101-tools/venv
fi

source orion101-tools/venv/bin/activate

find orion101-tools -name requirements.txt -exec cat {} \; -exec echo \; | sort -u > requirements.txt
uv pip install -r requirements.txt
