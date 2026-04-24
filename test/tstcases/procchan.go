package main

type ProcChan struct {
	Exitedchl chan int
	Exitchl   chan int
	Sendchl   chan string
	Rcvchl    chan string
	rcvlen    int
	sndlen    int
}

func NewProcChan(sndsize, rcvsize int) *ProcChan {
	var ptr *ProcChan
	ptr = &ProcChan{}
	ptr.Exitedchl = make(chan int, 10)
	ptr.Exitchl = make(chan int, 10)
	ptr.Sendchl = make(chan string, sndsize)
	ptr.Rcvchl = make(chan string, rcvsize)
	ptr.rcvlen = rcvsize
	ptr.sndlen = sndsize
	return ptr
}

func (ptr *ProcChan) GetSendSize() int {
	return ptr.sndlen
}

func (ptr *ProcChan) GetRcvSize() int {
	return ptr.rcvlen
}
