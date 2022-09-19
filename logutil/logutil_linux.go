package logutil

type unixbackGround struct {
}

func nativeGround() *unixbackGround {
	p := &unixbackGround{}
	return p
}

func (pb *unixbackGround) LogDebugOutputBackGround(s string) error {
	return nil
}

func (pb *unixbackGround) CloseDebugOutputBackGround() error {
	return nil
}
