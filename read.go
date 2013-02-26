// Package logplex implements streaming of syslog messages
package logplex

import (
	"io"
	"runtime"
	"time"
	"log"
)

type Msg struct {
	Priority  int
	Timestamp []byte
	Host      []byte
	User      []byte
	Pid       []byte
	Id        []byte
	Msg       []byte
}

func (m *Msg) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, string(m.Timestamp))
}

type BytesReader interface {
	io.Reader
	ReadLine() (line []byte, isPrefix bool, err error)
}

// Reader reads syslog streams
type Reader struct {
	buf BytesReader
}

// NewReader returns a new Reader that reads from buf.
func NewReader(buf BytesReader) *Reader {
	return &Reader{buf: buf}
}

// ReadMsg returns a single Msg. If no data is available, returns an error.
func (r *Reader) ReadMsg() (m *Msg, err error) {
	defer errRecover(&err)

	b := r.next()

	m = new(Msg)
	m.Priority = b.priority()
	m.Timestamp = b.bytes()
	m.Host = b.bytes()
	m.User = b.bytes()
	m.Pid = b.bytes()
	m.Id = b.bytes()
	m.Msg = b

	return
}

func (r *Reader) next() readBuf {
	var rbuf []byte
	buf, more, err := r.buf.ReadLine()
	if err != nil {
		log.Println(err)
		return rbuf
	}
	rbuf = append(rbuf, buf...)
	if more {
		for {
			tbuf, more, err := r.buf.ReadLine()
			if err != nil {
				log.Println(err)
				break
			}
			rbuf = append(rbuf, tbuf...)
			if !more {
				break
			}
		}
	}
	return rbuf
}

func errRecover(err *error) {
	e := recover()
	if e != nil {
		switch ee := e.(type) {
		case runtime.Error:
			panic(e)
		case error:
			*err = ee
		default:
			panic(e)
		}
	}
}
