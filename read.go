// Package logplex implements streaming of syslog messages
package logplex

import (
	"io"
	"runtime"
	"time"
)

type Msg struct {
	Priority     int
	Timestamp    []byte
	Host         []byte
	AppName      []byte
	Pid          []byte
	Id           []byte
	StructedData []byte
	Msg          []byte
}

func (m *Msg) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, string(m.Timestamp))
}

type BytesReader interface {
	io.Reader
	ReadString(delim byte) (line string, err error)
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

	b, e := r.next()
	if e != nil {
		return nil, e
	}
	m = new(Msg)
	m.Priority = b.priority()
	m.Timestamp = b.bytes()
	m.Host = b.bytes()
	m.AppName = b.bytes()
	m.Pid = b.bytes()
	m.Id = b.bytes()
	m.StructedData = b.bytes()
	m.Msg = b
	return
}

func (r *Reader) next() (readBuf, error) {
	ln, err := r.buf.ReadString('\n')
	return []byte(ln), err
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
