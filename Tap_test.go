package netaux_test

import (
	"bytes"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/maxbet1507/netaux"
)

type ctap struct {
	Called bool
}

func (s *ctap) Close() error {
	s.Called = true
	return nil
}

func TestTap(t *testing.T) {
	var conn1, conn2 net.Conn

	addr1 := "dummy1"
	addr2 := "dummy2"
	conn1, conn2 = netaux.FakePipe(addr1, addr2)

	rtap := &bytes.Buffer{}
	wtap := &bytes.Buffer{}
	ctap := &ctap{}
	tconn := netaux.Tap(conn1, rtap, wtap, ctap)

	if v := tconn.LocalAddr(); v != conn1.LocalAddr() {
		t.Fatal(v)
	}
	if v := tconn.RemoteAddr(); v != conn1.RemoteAddr() {
		t.Fatal(v)
	}
	if err := tconn.SetDeadline(time.Now()); err != conn1.SetDeadline(time.Now()) {
		t.Fatal(err)
	}
	if err := tconn.SetReadDeadline(time.Now()); err != conn1.SetReadDeadline(time.Now()) {
		t.Fatal(err)
	}
	if err := tconn.SetWriteDeadline(time.Now()); err != conn1.SetWriteDeadline(time.Now()) {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		n, err := tconn.Write([]byte("hello"))
		if n != 5 || err != nil {
			t.Fatal(n, err)
		}

		buf := make([]byte, 10)
		n, err = tconn.Read(buf)
		if n != 5 || err != nil || string(buf[:n]) != "world" {
			t.Fatal(n, err, buf)
		}

		if err := tconn.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()

		buf := make([]byte, 10)
		n, err := conn2.Read(buf)
		if n != 5 || err != nil || string(buf[:n]) != "hello" {
			t.Fatal(n, err, buf)
		}

		n, err = conn2.Write([]byte("world"))
		if n != 5 || err != nil {
			t.Fatal(n, err)
		}

		if err := conn2.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	wg.Wait()

	if v := wtap.String(); v != "hello" {
		t.Fatal(v)
	}
	if v := rtap.String(); v != "world" {
		t.Fatal(v)
	}
	if !ctap.Called {
		t.Fatal(ctap)
	}
}
