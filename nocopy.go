package stringx

type nocopy struct{}

func (*nocopy) Lock() {}

func (*nocopy) Unlock() {}
