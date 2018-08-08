package netaux

import (
	"fmt"
	"io"
	"net"
	"time"
)

var (
	errNotImplemented = fmt.Errorf("not implemented")
)

type fakeConn struct {
	Reader        io.Reader
	Writer        io.Writer
	Closer        io.Closer
	LocalAddress  net.Addr
	RemoteAddress net.Addr
}

func (s *fakeConn) LocalAddr() net.Addr {
	return s.LocalAddress
}

func (s *fakeConn) RemoteAddr() net.Addr {
	return s.RemoteAddress
}

func (s *fakeConn) SetDeadline(time.Time) error {
	return errNotImplemented
}

func (s *fakeConn) SetReadDeadline(time.Time) error {
	return errNotImplemented
}

func (s *fakeConn) SetWriteDeadline(time.Time) error {
	return errNotImplemented
}

func (s *fakeConn) Read(p []byte) (n int, err error) {
	err = errNotImplemented
	if s.Reader != nil {
		n, err = s.Reader.Read(p)
	}
	return
}

func (s *fakeConn) Write(p []byte) (n int, err error) {
	err = errNotImplemented
	if s.Writer != nil {
		n, err = s.Writer.Write(p)
	}
	return
}

func (s *fakeConn) Close() (err error) {
	err = errNotImplemented
	if s.Closer != nil {
		err = s.Closer.Close()
	}
	return
}

// FakeConn -
func FakeConn(l, r string, reader io.Reader, writer io.Writer, closer io.Closer) net.Conn {
	ret := &fakeConn{
		Reader:        reader,
		Writer:        writer,
		Closer:        closer,
		LocalAddress:  FakeAddr(l),
		RemoteAddress: FakeAddr(r),
	}
	return ret
}

// FakePipe -
func FakePipe(l, r string) (net.Conn, net.Conn) {
	oir, oiw := io.Pipe()
	ior, iow := io.Pipe()

	oc := FakeConn(l, r, oir, iow, iow)
	ic := FakeConn(r, l, ior, oiw, oiw)

	return oc, ic
}
