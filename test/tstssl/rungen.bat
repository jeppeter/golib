rem ECHO off

set CURDIR=%~dp0

%CURDIR%\tstssl.exe rsagen --keyfile %CURDIR%\root.key 4096
%CURDIR%\tstssl.exe rsagen --keyfile %CURDIR%\intern.key 4096
%CURDIR%\tstssl.exe rsagen --keyfile %CURDIR%\intern2.key 4096
%CURDIR%\tstssl.exe rsagen --keyfile %CURDIR%\child.key 4096

%CURDIR%\tstssl.exe x509create --keyfile %CURDIR%\root.key %CURDIR%\root.json --certfile %CURDIR%\root.pem
%CURDIR%\tstssl.exe x509create --keyfile %CURDIR%\intern.key %CURDIR%\intern.json %CURDIR%\root.pem %CURDIR%\root.key --certfile %CURDIR%\intern.pem
%CURDIR%\tstssl.exe x509create --keyfile %CURDIR%\intern2.key %CURDIR%\intern2.json %CURDIR%\root.pem %CURDIR%\root.key --certfile %CURDIR%\intern2.pem
%CURDIR%\tstssl.exe x509create --keyfile %CURDIR%\child.key %CURDIR%\child.json %CURDIR%\intern.pem %CURDIR%\intern.key --certfile %CURDIR%\child.pem

