#!/usr/bin/env bash

npm install

TARGET=${1:-dist}
STAGE=${2:-dev}

case "$STAGE" in
    "dev" )     ENV="development"   ;;
    "staging" ) ENV="staging"       ;;
    "prod" )    ENV="production"    ;;
esac

sls deploy --env $ENV --stage $STAGE --config ./deployments/serverless_$TARGET.yml
