echo off
set BATDIR=%~dp0

copy /y %BATDIR%\..\..\log\loglib.go .
copy /y %BATDIR%\..\..\log\loglib_windows.go .
go build -o decode.exe decode.go loglib.go loglib_windows.go pkcs12.go
echo on