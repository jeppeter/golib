
set CURDIR=%~dp0
set TOPDIR=%CURDIR%\..\..

set GOPATH=%GOPATH%;%TOPDIR%
set GOOS=windows
set GOARCH=amd64
go build -o tstcases.exe main.go  mem_windows.go
