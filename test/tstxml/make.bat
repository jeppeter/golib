echo off
rem rmdir /s /q src || echo ""
set CURDIR=%~dp0
set TOPDIR=%CURDIR%\..\..

REM md src\jsonext
REM xcopy /s /e %CD%\..\..\jsonext src\jsonext
set GOPATH=%GOPATH%;%TOPDIR%
go build -o tstxml.exe main.go