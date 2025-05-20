#!/bin/bash

# Create build directory
mkdir -p build

# Build the Go application
go build -o build/tt cmd/*.go

if [ -d build/test ]
then
    cp build/tt build/test/tt
fi