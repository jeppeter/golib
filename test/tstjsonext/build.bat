echo off
rmdir /s /q src || echo ""
set CURDIR=%~dp0
set TOPDIR=%CURDIR%\..\..

REM md src\jsonext
REM xcopy /s /e %CD%\..\..\jsonext src\jsonext
set GOPATH=%GOPATH%;%TOPDIR%
go build -o tstjsonext.exe main.go