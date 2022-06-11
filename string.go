package st

type String struct {
	mem []byte
	len int
	cap int
}

func (s *String) grow(n int) {
	s.mem = append(s.mem, make([]byte, n)...)
	s.cap += n
}

func (s *String) payload() []byte {
	return s.mem[0:s.len]
}

func (s *String) block() []byte {
	return s.mem
}
