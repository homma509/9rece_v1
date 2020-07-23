#!/usr/bin/env bash

npm install

IO=${1:-func}
STAGE=${2:-dev}

case "$STAGE" in
    "dev" )     ENV="development"   ;;
    "staging" ) ENV="staging"       ;;
    "prod" )    ENV="production"    ;;
esac

sls deploy --env $ENV --config ./deployments/serverless_$IO.yml --stage $STAGE
