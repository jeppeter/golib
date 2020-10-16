rmdir /s /q src || echo ""
md src\jsonext
xcopy /s /e %CD%\..\..\jsonext src\jsonext
set GOPATH=%GOPATH%;%CD%
go build -o tstjsonext.exe main.go