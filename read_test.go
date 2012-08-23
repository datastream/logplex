package logplex

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	data := strings.NewReader(
		"67 <174>1 2012-07-22T00:06:26-00:00 somehost Go console 2 Hi from Go\n" +
			"67 <174>1 2012-07-22T00:06:26-00:00 somehost Go console 10 Hi from Go",
	)

	exp := []*Msg{
		{
			174,
			[]byte("2012-07-22T00:06:26-00:00"),
			[]byte("somehost"),
			[]byte("Go"),
			[]byte("console"),
			[]byte("2"),
			[]byte("Hi from Go\n"),
		},
		{
			174,
			[]byte("2012-07-22T00:06:26-00:00"),
			[]byte("somehost"),
			[]byte("Go"),
			[]byte("console"),
			[]byte("10"),
			[]byte("Hi from Go"),
		},
	}

	r := NewReader(bufio.NewReader(data))

	for i, e := range exp {
		m, err := r.ReadMsg()
		if err != nil {
			t.Errorf("error on %d: %v", i, err)
			continue
		}
		if m.Priority != e.Priority {
			t.Errorf("expected %d, got %d", m.Priority, e.Priority)
		}
		if !reflect.DeepEqual(m.Timestamp, e.Timestamp) {
			t.Errorf("expected %d, got %d", m.Timestamp, e.Timestamp)
		}
		if !reflect.DeepEqual(m.Host, e.Host) {
			t.Errorf("expected %s, got %s", m.Host, e.Host)
		}
		if !reflect.DeepEqual(m.Pid, e.Pid) {
			t.Errorf("expected %s, got %s", m.Pid, e.Pid)
		}
		if !reflect.DeepEqual(m.Id, e.Id) {
			t.Errorf("expected %s, got %s", m.Id, e.Id)
		}
		if !reflect.DeepEqual(m.Msg, e.Msg) {
			t.Errorf("expected %s, got %s", m.Msg, e.Msg)
		}
	}

	_, err := r.ReadMsg()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}

	// one more for good measure
	_, err = r.ReadMsg()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}
