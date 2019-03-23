package newrelic

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"testing"
)

type rwNoExtraMethods struct {
	hijackCalled      bool
	readFromCalled    bool
	flushCalled       bool
	closeNotifyCalled bool
}

type rwTwoExtraMethods struct{ rwNoExtraMethods }
type rwAllExtraMethods struct{ rwTwoExtraMethods }

func (rw *rwAllExtraMethods) CloseNotify() <-chan bool {
	rw.closeNotifyCalled = true
	return nil
}
func (rw *rwAllExtraMethods) ReadFrom(r io.Reader) (int64, error) {
	rw.readFromCalled = true
	return 0, nil
}

func (rw *rwNoExtraMethods) Header() http.Header        { return nil }
func (rw *rwNoExtraMethods) Write([]byte) (int, error)  { return 0, nil }
func (rw *rwNoExtraMethods) WriteHeader(statusCode int) {}

func (rw *rwTwoExtraMethods) Flush() {
	rw.flushCalled = true
}
func (rw *rwTwoExtraMethods) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw.hijackCalled = true
	return nil, nil, nil
}

func TestTransactionAllExtraMethods(t *testing.T) {
	app := testApp(nil, nil, t)
	rw := &rwAllExtraMethods{}
	txn := app.StartTransaction("hello", rw, nil)
	if v, ok := txn.(http.CloseNotifier); ok {
		v.CloseNotify()
	}
	if v, ok := txn.(http.Flusher); ok {
		v.Flush()
	}
	if v, ok := txn.(http.Hijacker); ok {
		v.Hijack()
	}
	if v, ok := txn.(io.ReaderFrom); ok {
		v.ReadFrom(nil)
	}
	if !rw.hijackCalled ||
		!rw.readFromCalled ||
		!rw.flushCalled ||
		!rw.closeNotifyCalled {
		t.Error("wrong methods called", rw)
	}
}

func TestTransactionNoExtraMethods(t *testing.T) {
	app := testApp(nil, nil, t)
	rw := &rwNoExtraMethods{}
	txn := app.StartTransaction("hello", rw, nil)
	if _, ok := txn.(http.CloseNotifier); ok {
		t.Error("unexpected CloseNotifier method")
	}
	if _, ok := txn.(http.Flusher); ok {
		t.Error("unexpected Flusher method")
	}
	if _, ok := txn.(http.Hijacker); ok {
		t.Error("unexpected Hijacker method")
	}
	if _, ok := txn.(io.ReaderFrom); ok {
		t.Error("unexpected ReaderFrom method")
	}
}

func TestTransactionTwoExtraMethods(t *testing.T) {
	app := testApp(nil, nil, t)
	rw := &rwTwoExtraMethods{}
	txn := app.StartTransaction("hello", rw, nil)
	if _, ok := txn.(http.CloseNotifier); ok {
		t.Error("unexpected CloseNotifier method")
	}
	if v, ok := txn.(http.Flusher); ok {
		v.Flush()
	}
	if v, ok := txn.(http.Hijacker); ok {
		v.Hijack()
	}
	if _, ok := txn.(io.ReaderFrom); ok {
		t.Error("unexpected ReaderFrom method")
	}
	if !rw.hijackCalled ||
		rw.readFromCalled ||
		!rw.flushCalled ||
		rw.closeNotifyCalled {
		t.Error("wrong methods called", rw)
	}
}
