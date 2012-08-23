package logplex

import (
	"bufio"
	"io"
	"strconv"
	"time"
	"runtime"
	"errors"
)

var (
	ErrInvalidPriority = errors.New("Invalid Priority")
)

type Msg struct {
	Priority int
	Time     time.Time
	Host     []byte
	User     []byte
	Pid      []byte
	Id       []byte
	Msg      []byte
}

type Reader struct {
	buf  *bufio.Reader
	err  error
}

func NewReader(buf *bufio.Reader) *Reader {
	return &Reader{buf: buf}
}

func (r *Reader) ReadMsg() (m *Msg, err error) {
	defer errRecover(&err)

	b := r.next()

	m = new(Msg)
	m.Priority = b.priority()
	m.Time = b.time()
	m.Host = b.bytes()
	m.User = b.bytes()
	m.Pid = b.bytes()
	m.Id = b.bytes()
	m.Msg = b

	return
}

func (r *Reader) next() readBuf {
	if r.err != nil {
		return nil
	}

	b, err := r.buf.ReadBytes(' ')
	if err != nil {
		panic(err)
	}
	b = b[:len(b)-1]

	n, err := strconv.Atoi(string(b))
	if err != nil {
		panic(err)
	}

	buf := make(readBuf, n-1)
	_, err = io.ReadFull(r.buf, buf)
	if err != nil {
		panic(err)
	}

	return buf
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