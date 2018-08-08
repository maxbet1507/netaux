package netaux

// FakeAddr -
type FakeAddr string

// Network -
func (s FakeAddr) Network() string {
	return string(s)
}

func (s FakeAddr) String() string {
	return string(s)
}
