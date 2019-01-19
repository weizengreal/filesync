#/bin/bash

export GOPATH=$GOPATH:`pwd`

go build -i -o ./target/filesync ./src/main.go

./target/filesync --help
