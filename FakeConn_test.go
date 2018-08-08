package netaux_test

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/maxbet1507/netaux"
)

func TestFakePipe(t *testing.T) {
	var conn1, conn2 net.Conn

	addr1 := "dummy1"
	addr2 := "dummy2"
	conn1, conn2 = netaux.FakePipe(addr1, addr2)

	if v := conn1.LocalAddr(); v.String() != addr1 {
		t.Fatal(v)
	}
	if v := conn2.LocalAddr(); v.String() != addr2 {
		t.Fatal(v)
	}

	if v := conn1.RemoteAddr(); v.String() != addr2 {
		t.Fatal(v)
	}
	if v := conn2.RemoteAddr(); v.String() != addr1 {
		t.Fatal(v)
	}

	if err := conn1.SetDeadline(time.Now()); err == nil {
		t.Fatal(err)
	}
	if err := conn2.SetDeadline(time.Now()); err == nil {
		t.Fatal(err)
	}

	if err := conn1.SetReadDeadline(time.Now()); err == nil {
		t.Fatal(err)
	}
	if err := conn2.SetReadDeadline(time.Now()); err == nil {
		t.Fatal(err)
	}

	if err := conn1.SetWriteDeadline(time.Now()); err == nil {
		t.Fatal(err)
	}
	if err := conn2.SetWriteDeadline(time.Now()); err == nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		n, err := conn1.Write([]byte("hello"))
		if n != 5 || err != nil {
			t.Fatal(n, err)
		}

		buf := make([]byte, 10)
		n, err = conn1.Read(buf)
		if n != 5 || err != nil || string(buf[:n]) != "world" {
			t.Fatal(n, err, buf)
		}

		if err := conn1.Close(); err != nil {
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
}
