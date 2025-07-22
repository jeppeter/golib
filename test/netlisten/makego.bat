echo off
set CURDIR=%~dp0
set TOPDIR=%CURDIR%\..\..
set GOOS=windows
set GOPATH=%GOPATH%;%TOPDIR%
go build -o nc.exe nc.go

