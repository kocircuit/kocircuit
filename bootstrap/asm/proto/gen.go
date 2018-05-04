package proto

//go:generate protoc -I=. -I=$GOPATH/src --gofast_out=. asm.proto

func FileDescriptorBytes() []byte {
	return fileDescriptorAsm
}
