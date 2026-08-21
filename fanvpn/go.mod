module fanvpn

go 1.25

replace dbgutil => ../src/dbgutil

replace fileop => ../src/fileop

replace aesext => ../src/aesext

replace logutil => ../src/logutil

require (
	aesext v0.0.0-00010101000000-000000000000 // indirect
	dbgutil v0.0.0-00010101000000-000000000000 // indirect
	fileop v0.0.0-00010101000000-000000000000 // indirect
	github.com/jeppeter/go-extargsparse v1.0.0 // indirect
	github.com/jeppeter/log4go v0.0.0-20191224035337-ce096513aa1f // indirect
	github.com/tebeka/atexit v0.3.0 // indirect
	logutil v0.0.0-00010101000000-000000000000 // indirect
)
