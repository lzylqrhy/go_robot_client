package protocol

type C2SPing struct {
	TimeStamp uint32
}

func (p *C2SPing) Bytes() []byte {
	var pb Protocol
	pb.SetCmd(0x1)
	return pb.Bytes()
}
