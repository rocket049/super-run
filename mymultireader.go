package main

import (
	"io"
	"sync/atomic"
)

type RandMultiReader struct {
	channel chan []byte
	buf     []byte
	num     int32
}

func (s *RandMultiReader) Read(p []byte) (n int, err error) {
	lbuf := len(s.buf)
	lp := len(p)
	if lbuf > 0 {
		if lbuf <= lp {
			n = copy(p, s.buf)
			s.buf = nil
		} else {
			n = copy(p, s.buf[:lp])
			s.buf = s.buf[lp:]
		}
	} else {
		var ok bool
		s.buf, ok = <-s.channel
		//log.Println("Read ", ok)
		if !ok {
			return 0, io.EOF
		} else {
			lbuf = len(s.buf)
			if lbuf <= lp {
				n = copy(p, s.buf)
				s.buf = nil
			} else {
				n = copy(p, s.buf[:lp])
				s.buf = s.buf[lp:]
			}
		}
	}
	//log.Printf("Read:%d\n", n)
	return
}

func (s *RandMultiReader) LinkReader(r io.ReadCloser) error {
	defer r.Close()
	for {
		p := make([]byte, 512)
		n, err := r.Read(p)
		if err != nil {
			atomic.AddInt32(&s.num, -1)
			if atomic.LoadInt32(&s.num) == 0 {
				close(s.channel)
				//log.Println("close channel")
			}
			return err
		}
		s.channel <- p[:n]
	}

	return nil
}

func NewRandMultiReader(readers ...io.ReadCloser) io.Reader {
	res := &RandMultiReader{channel: make(chan []byte, 10), buf: nil, num: int32(len(readers))}
	for _, v := range readers {
		go res.LinkReader(v)
	}
	return res
}
