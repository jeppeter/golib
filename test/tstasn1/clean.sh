#! /bin/bash

scriptfile=`readlink -f $0`
scriptdir=`dirname $scriptfile`

rm -f $scriptdir/golib.go $scriptdir/golib_unix.go $scriptdir/decode