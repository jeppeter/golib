rem echo off

set CURDIR=%~dp0
%CURDIR%\tstssl.exe x509vfy --vfyopt %CURDIR%\verify.json %CURDIR%\child.pem -vvvvv