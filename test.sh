#!/bin/bash
export GOPATH="$(pwd)"/code
go test $1 -test.v -test.parallel 2
