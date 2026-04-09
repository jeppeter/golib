module tstsockproto

go 1.25.0

replace dbgutil => ../../src/dbgutil

replace logutil => ../../src/logutil

replace socktimeout => ../../src/socktimeout

replace sockproto => ../../src/sockproto

replace jsonproto => ../../src/jsonproto

replace fileop => ../../src/fileop

require (
	dbgutil v0.0.0-00010101000000-000000000000 // indirect
	fileop v0.0.0-00010101000000-000000000000 // indirect
	github.com/jeppeter/go-extargsparse v1.0.0 // indirect
	github.com/jeppeter/log4go v0.0.0-20191224035337-ce096513aa1f // indirect
	github.com/tebeka/atexit v0.3.0 // indirect
	jsonproto v0.0.0-00010101000000-000000000000 // indirect
	logutil v0.0.0-00010101000000-000000000000 // indirect
	sockproto v0.0.0-00010101000000-000000000000 // indirect
	socktimeout v0.0.0-00010101000000-000000000000 // indirect
)
