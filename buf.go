package logplex

import (
	"bytes"
	"io"
	"strconv"
	"time"
)

type readBuf []byte

func (b *readBuf) int() int {
	p := b.bytes()
	n, err := strconv.Atoi(string(p))
	if err != nil {
		panic(err)
	}
	*b = (*b)[len(p):]
	return n
}

func (b *readBuf) bytes() []byte {
	i := bytes.IndexByte(*b, ' ')
	if i < 0 {
		panic(io.ErrUnexpectedEOF)
	}
	bs := (*b)[:i]
	*b = (*b)[i+1:]
	return bs
}

func (b *readBuf) priority() int {
	p := b.bytes()
	if len(p) < 4 {
		panic(ErrInvalidPriority)
	}
	n, err := strconv.Atoi(string(p[1 : len(p)-2]))
	if err != nil {
		panic(ErrInvalidPriority)
	}
	return n
}

func (b *readBuf) time() time.Time {
	return mustParseTime(time.RFC3339, string(b.bytes()))
}

func mustParseTime(format string, s string) time.Time {
	t, err := time.Parse(format, s)
	if err != nil {
		panic(err)
	}
	return t
}
