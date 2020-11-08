#!/usr/bin/env bash

cd api
find . -name main.go -type f \
 | xargs -n 1 dirname \
 | xargs -n 1 -I@ bash -c "CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o ../build/@/main ./@"