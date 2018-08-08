package netaux

import (
	"bytes"
	"io"
	"net"
	"time"
)

type tapConn struct {
	Conn     net.Conn
	TapRead  io.Writer
	TapWrite io.Writer
	TapClose []io.Closer
}

func (s *tapConn) LocalAddr() net.Addr {
	return s.Conn.LocalAddr()
}

func (s *tapConn) RemoteAddr() net.Addr {
	return s.Conn.RemoteAddr()
}

func (s *tapConn) SetDeadline(t time.Time) error {
	return s.Conn.SetDeadline(t)
}

func (s *tapConn) SetReadDeadline(t time.Time) error {
	return s.Conn.SetReadDeadline(t)
}

func (s *tapConn) SetWriteDeadline(t time.Time) error {
	return s.Conn.SetWriteDeadline(t)
}

func (s *tapConn) Read(p []byte) (int, error) {
	n, err := s.Conn.Read(p)
	if n > 0 && s.TapWrite != nil {
		io.Copy(s.TapRead, bytes.NewBuffer(p[:n])) // ignore errors
	}
	return n, err
}

func (s *tapConn) Write(p []byte) (int, error) {
	n, err := s.Conn.Write(p)
	if n > 0 && s.TapWrite != nil {
		io.Copy(s.TapWrite, bytes.NewBuffer(p[:n])) // ignore errors
	}
	return n, err
}

func (s *tapConn) Close() error {
	for _, c := range s.TapClose {
		if c != nil {
			c.Close()
		}
	}
	return s.Conn.Close()
}

// Tap -
func Tap(conn net.Conn, r, w io.Writer, c ...io.Closer) net.Conn {
	ret := &tapConn{
		Conn:     conn,
		TapRead:  r,
		TapWrite: w,
		TapClose: c,
	}
	return ret
}
