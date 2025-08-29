#!/bin/bash

set -xe

mkdir -p bin
go build -o bin/main src/main.go
