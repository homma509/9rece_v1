#!/usr/bin/env bash

if [ -d .serverless ]; then
	rm -R .serverless
fi
if [ -d build ]; then
	rm -R build
fi
if [ -d test/cover ]; then
	rm -R test/cover
fi