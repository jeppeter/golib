
set CURDIR=%~dp0
copy /y %CURDIR%..\..\log\loglib.go loglib.go
copy /y %CURDIR%..\..\log\loglib_windows.go loglib_windows.go
copy /y %CURDIR%..\..\pipecmdlib\pipecmdlib_windows.go pipecmdlib_windows.go
copy /y %CURDIR%..\..\pipecmdlib\pipecmdlib.go pipecmdlib.go
go build -o pipecmd.exe pipecmdlib.go pipecmdlib_windows.go main.go loglib.go loglib_windows.go