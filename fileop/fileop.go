package main

func read_file_bytes(fname string) (rbytes []byte, err error) {
	err = nil
	rbytes = []byte{}
	return
}

func write_file_bytes(fname string, obytes []byte) (err error) {
	err = nil
	return
}

func read_file(fname string) (s string, err error) {
	var ob []byte
	s = ""
	ob, err = read_file_bytes(fname)
	if err != nil {
		return
	}
	s = string(ob)
	return
}

func write_file(fname string, ostring string) (err error) {
	err = write_file_bytes(fname, []byte(ostring))
	return
}
