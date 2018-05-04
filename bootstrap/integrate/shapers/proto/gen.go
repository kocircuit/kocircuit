package proto

//go:generate protoc -I=. -I=$GOPATH/src --gofast_out=. shapers.proto

func FileDescriptorBytes() []byte {
	return fileDescriptorShapers
}
