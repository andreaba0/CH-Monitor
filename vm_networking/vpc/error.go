package networkvpc

type ErrNotEnoughBytes struct{}

func (err *ErrNotEnoughBytes) Error() string {
	return "not enough bytes to process the packet"
}
