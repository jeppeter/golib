REM to copy files
SET scrdir=%~dp0
copy /y %scrdir%..\log\loglib.go %scrdir%loglib.go
copy /y %scrdir%..\log\loglib_windows.go %scrdir%loglib_windows.go
go build -o httpsvr.exe main.go loglib.go loglib_windows.go server.go