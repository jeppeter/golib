module stockhdl

go 1.25.0

replace logutil => ../src/logutil

replace dbgutil => ../src/dbgutil

replace fileop => ../src/fileop

replace github.com/jeppeter/go-sqlite3dyn => ../../go-sqlite3dyn

require (
	dbgutil v0.0.0-00010101000000-000000000000 // indirect
	fileop v0.0.0-00010101000000-000000000000 // indirect
	github.com/jeppeter/go-extargsparse v1.0.0 // indirect
	github.com/jeppeter/go-sqlite3dyn v1.0.0 // indirect
	github.com/jeppeter/log4go v0.0.0-20191224035337-ce096513aa1f // indirect
	github.com/notti/nocgo v0.0.0-20190619201224-fc443047424c // indirect
	github.com/tebeka/atexit v0.3.0 // indirect
	logutil v0.0.0-00010101000000-000000000000 // indirect
)
