
set GOPATH=%GOPATH%;%CD%\..\..\
set GOOS=windows
go build -o tstcases.exe main.go 
set GOOS=linux
go build -o tstcases main.go 