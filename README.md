# Get Start

```shell
git clone https://github/fudali113/mongosync mongosync/src/github/fudali113/mongosync
export GOPATH=$pwd/mongosync:$GOPATH
cd mongosync/src/github/fudali113/mongosync
go run main.go
```


get start:
```
./mongosync -help
./mongosync -src=192.168.30.249:28717 -dst=localhost:27017 -name=test
```