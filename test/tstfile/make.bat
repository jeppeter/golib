echo off
set CURDIR=%~dp0
set TOPDIR=%CURDIR%\..\..
set GOOS=windows
set GOPATH=%GOPATH%;%TOPDIR%
go build -o tstfile.exe main.go

