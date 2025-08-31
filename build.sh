#!/bin/bash

set -xe

rm -r bin
mkdir -p bin
go build -o bin/main src/main.go
