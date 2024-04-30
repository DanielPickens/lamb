#!/bin/bash

set -e
# set the working directory to the location of this script for building the docker image through pkger and makefile
go get github.com/markbates/pkger/cmd/pkger
make build
docker cp ./ e2e-command-runner:/lamb
