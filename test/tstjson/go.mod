module tstjson

go 1.25.0

replace dbgutil => ../../src/dbgutil

replace logutil => ../../src/logutil

replace fileop => ../../src/fileop

replace jsonext => ../../src/jsonext

require (
	dbgutil v0.0.0-00010101000000-000000000000
	fileop v0.0.0-00010101000000-000000000000
	github.com/jeppeter/go-extargsparse v1.0.0
	github.com/tebeka/atexit v0.3.0
	jsonext v0.0.0-00010101000000-000000000000
	logutil v0.0.0-00010101000000-000000000000
)

require github.com/jeppeter/log4go v0.0.0-20191224035337-ce096513aa1f // indirect
