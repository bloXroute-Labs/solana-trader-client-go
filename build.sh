#!/usr/bin/env bash

IMAGE=bloxroute/solana-trader-client-go:${1:-latest}
echo "building container... $IMAGE"
echo ""

docker build . -f Dockerfile --rm=true -t $IMAGE --platform linux/amd64
