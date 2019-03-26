#!/bin/bash -e
path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        goimports -w $path
        gofmt -w $path
        ;;
esac
go test -coverprofile /tmp/c.out -args -update-golden
uncover /tmp/c.out
