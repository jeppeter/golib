echo off
rem rmdir /s /q src || echo ""
set CURDIR=%~dp0
set TOPDIR=%CURDIR%\..\..

REM md src\jsonext
REM xcopy /s /e %CD%\..\..\jsonext src\jsonext
set GOPATH=%GOPATH%;%TOPDIR%
go build -o tstssl.exe main.go x509temp.go x509parse.go x509disp.go x509create.go defines.go x509vfy.go jsonmap.go ncert.go tstrsa.go pemlib.go