package netaux_test

import (
	"net"
	"testing"

	"github.com/maxbet1507/netaux"
)

func TestFakeAddr(t *testing.T) {
	var addr net.Addr

	addr = netaux.FakeAddr("dummy")
	if v := addr.Network(); v != "dummy" {
		t.Fatal(v)
	}
	if v := addr.String(); v != "dummy" {
		t.Fatal(v)
	}
}
