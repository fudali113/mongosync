#!/bin/sh
go build -o bin/mongosync
env GOOS=linux GOARCH=arm64 go build -o bin/mongosync-liunx-64
env GOOS=darwin go build -o bin/mongosync-darwin
env GOOS=windows GOARCH=amd64 go build -o bin/mongosync-windows-64

echo build successed