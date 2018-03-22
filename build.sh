#!/bin/sh
go build -o bin/mongosync
env GOOS=linux GOARCH=amd64 go build -o bin/mongosync-linux-amd64
env GOOS=darwin go build -o bin/mongosync-darwin
env GOOS=windows GOARCH=amd64 go build -o bin/mongosync-windows-amd64

echo build successed