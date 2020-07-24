#!/usr/bin/env bash

npm install

STAGE=${1:-dev}

case "$STAGE" in
    "dev" )     ENV="development"   ;;
    "staging" ) ENV="staging"       ;;
    "prod" )    ENV="production"    ;;
esac

sls remove --env $ENV --config ./deployments/serverless.yml --stage $STAGE
