module tstfile

go 1.25.0

replace logutil => ../../src/logutil

replace fileop => ../../src/fileop

replace dbgutil => ../../src/dbgutil

require (
	fileop v0.0.0-00010101000000-000000000000
	github.com/jeppeter/go-extargsparse v1.0.0
	github.com/tebeka/atexit v0.3.0
	logutil v0.0.0-00010101000000-000000000000
)

require (
	dbgutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/jeppeter/log4go v0.0.0-20191224035337-ce096513aa1f // indirect
	golang.org/x/net v0.58.0 // indirect
)
