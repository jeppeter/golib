echo off
set CURDIR=%~dp0
set TOPDIR=%CURDIR%\..\..
set GOOS=windows
set GOPATH=%GOPATH%;%TOPDIR%
go build -o tstexecutil.exe main.go

set GOOS=linux
go build -o tstexecutil main.go
