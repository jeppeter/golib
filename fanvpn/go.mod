module fanvpn

go 1.25.0

replace dbgutil => ../src/dbgutil

replace fileop => ../src/fileop

replace aesext => ../src/aesext

replace logutil => ../src/logutil

replace tcprelay => ../src/tcprelay

replace socktimeout => ../src/socktimeout

require (
	aesext v0.0.0-00010101000000-000000000000
	dbgutil v0.0.0-00010101000000-000000000000
	fileop v0.0.0-00010101000000-000000000000
	github.com/jeppeter/go-extargsparse v1.0.0
	github.com/tebeka/atexit v0.3.0
	logutil v0.0.0-00010101000000-000000000000
	tcprelay v0.0.0-00010101000000-000000000000
)

require (
	github.com/jeppeter/log4go v0.0.0-20191224035337-ce096513aa1f // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	socktimeout v0.0.0-00010101000000-000000000000 // indirect
)
