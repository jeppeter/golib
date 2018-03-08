#! /bin/bash

scriptfile=`readlink -f $0`
scriptdir=`dirname $scriptfile`

cp -f $scriptdir/../../log/loglib.go .
cp -f $scriptdir/../../log/loglib_unix.go .
go build -o main main.go loglib.go loglib_unix.go
