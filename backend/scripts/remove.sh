#!/usr/bin/env bash

npm install

TARGET=${1:-func}
STAGE=${2:-dev}

case "$STAGE" in
    "dev" )     ENV="development"   ;;
    "staging" ) ENV="staging"       ;;
    "prod" )    ENV="production"    ;;
esac

sls remove --env $ENV --config ./deployments/serverless_$TARGET.yml --stage $STAGE
