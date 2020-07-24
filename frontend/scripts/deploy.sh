#!/usr/bin/env bash

npm install

STAGE=${1:-dev}

case "$STAGE" in
    "dev" )     ENV="development"   ;;
    "staging" ) ENV="staging"       ;;
    "prod" )    ENV="production"    ;;
esac

sls deploy --env $ENV --stage $STAGE --config ./deployments/serverless.yml
