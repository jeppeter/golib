module tstcases

go 1.25.0

require github.com/jeppeter/go-extargsparse v1.0.0

require (
	dbgutil v0.0.0-00010101000000-000000000000 // indirect
	executil v0.0.0-00010101000000-000000000000 // indirect
	fileop v0.0.0-00010101000000-000000000000 // indirect
	github.com/jeppeter/log4go v0.0.0-20191224035337-ce096513aa1f // indirect
	github.com/jeppeter/npipe v0.0.0-20250421082957-07dee130f587 // indirect
	github.com/tebeka/atexit v0.3.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	logutil v0.0.0-00010101000000-000000000000 // indirect
	npipepack v0.0.0-00010101000000-000000000000 // indirect
	winpriv v0.0.0-00010101000000-000000000000 // indirect
	winreg v0.0.0-00010101000000-000000000000 // indirect
)

replace dbgutil => ../../src/dbgutil

replace executil => ../../src/executil

replace logutil => ../../src/logutil

replace fileop => ../../src/fileop

replace npipepack => ../../src/npipepack

replace winpriv => ../../src/winpriv

replace winreg => ../../src/winreg
