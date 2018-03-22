#!bin/sh
git clone https://git.23cube.com/tools/mongosync mongosync/src/github/fudali113/mongosync

export GOPATH=$pwd/mongosync:$GOPATH
cd mongosync/src/github/fudali113/mongosync