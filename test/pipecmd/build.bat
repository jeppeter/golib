
set CURDIR=%~dp0
copy /y %CURDIR%..\..\log\loglib.go loglib.go
copy /y %CURDIR%..\..\log\loglib_windows.go loglib_windows.go
go build -o pipecmd.exe pipecmd.go pipecmd_windows.go main.go loglib.go loglib_windows.go