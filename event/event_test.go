package event

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"
)

func TestWriteEscape(t *testing.T) {
	e := New(ioutil.Discard, 0, "", "")
	defer e.Close()
	e.Write([]byte("\b\f\r\n\t\\\""))
	if got, want := e.buf.String(), `\b\f\r\n\t\\\"`; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteMaxLen(t *testing.T) {
	e := New(ioutil.Discard, 5, "", "")
	defer e.Close()
	n, _ := e.Write([]byte("abcdefghij"))
	if got, want := n, 5; got != want {
		t.Errorf("invalid n: got %v, want %v", got, want)
	}
	if got, want := e.buf.String(), "abcde"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteMaxLenEscaping(t *testing.T) {
	e := New(ioutil.Discard, 5, "", "")
	defer e.Close()
	n, _ := e.Write([]byte("\n\n\n\n"))
	if got, want := n, 4; got != want {
		t.Errorf("invalid n: got %v, want %v", got, want)
	}
	if got, want := e.buf.String(), `\n\n`; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFlushEol(t *testing.T) {
	out := &bytes.Buffer{}
	e := New(out, 0, "\n", "")
	defer e.Close()
	e.Write([]byte("line1\n"))
	e.Write([]byte("line2"))
	e.Flush()
	if got, want := out.String(), "line1\\nline2\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFlushJSON(t *testing.T) {
	out := &bytes.Buffer{}
	e := New(out, 0, "\n", "message")
	defer e.Close()
	e.Write([]byte("line1\n"))
	e.Write([]byte("line2"))
	e.Flush()
	if got, want := out.String(), "{\"message\":\"line1\\nline2\"}\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFlushJSONMaxLen(t *testing.T) {
	out := &bytes.Buffer{}
	e := New(out, 25, "\n", "message")
	defer e.Close()
	e.Write([]byte("line1\n"))
	e.Write([]byte("line2"))
	e.Flush()
	if got, want := len(out.String()), 25; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := out.String(), "{\"message\":\"line1\\nlin\"}\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFlushEmpty(t *testing.T) {
	out := &bytes.Buffer{}
	e := New(out, 0, "\n", "")
	defer e.Close()
	e.Flush()
	if got, want := out.String(), ""; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEmpty(t *testing.T) {
	e := New(ioutil.Discard, 0, "\n", "")
	defer e.Close()
	if got, want := e.Empty(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	e.Write([]byte{' '})
	if got, want := e.Empty(), false; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	e.Flush()
	if got, want := e.Empty(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAutoFlush(t *testing.T) {
	done := make(chan bool, 1)
	autoFlushCalledHook = func() {
		done <- true
	}
	defer func() {
		autoFlushCalledHook = func() {}
	}()
	out := &bytes.Buffer{}
	e := New(out, 0, "\n", "")
	defer e.Close()
	e.Write([]byte{'x'})
	c := make(chan time.Time)
	e.start <- c // simulate AutoFlush()
	if got, want := out.String(), ""; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	go func() {
		c <- time.Time{}
	}()
	<-done
	if got, want := out.String(), "x\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
