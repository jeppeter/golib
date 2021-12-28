copy /y ..\..\log\loglib.go .
copy /y ..\..\log\loglib_windows.go .
copy /y ..\..\fileop\fileop.go .

go build -o tstcases.exe main.go loglib.go loglib_windows.go fileop.go