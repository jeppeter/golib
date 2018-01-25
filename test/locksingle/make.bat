@echo off

set BATDIR=%~dp0

copy /Y %BATDIR%\..\..\log\loglib.go .
copy /Y %BATDIR%\..\..\log\loglib_windows.go .
copy /Y %BATDIR%\..\..\locksingle\locksingle_windows.go .

go build -o locksingle.exe main.go loglib.go loglib_windows.go  locksingle_windows.go

echo on