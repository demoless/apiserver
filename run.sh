#!/bin/bash
cd $GOPATH/src/apiserver
gofmt -w .   
go tool vet .
go build -v .