#! /bin/bash

scriptfile=`readlink -f $0`
scriptdir=`dirname $scriptfile`

cp $scriptdir/../../log/loglib.go .
cp $scriptdir/../../log/loglib_unix.go .
cp $scriptdir/../../locksingle/locksingle_unix.go .

go build -o locksingle main.go loglib.go loglib_unix.go locksingle_unix.go
