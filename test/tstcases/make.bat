
set GOPATH=%GOPATH%;%CD%\..\..\
set GOOS=windows
set GOARCH=amd64
go build -o tstcases.exe main.go  mem_windows.go
set GOOS=linux
set GOARCH=amd64
go build -o tstcases main.go  mem_linux.go