package proto

//go:generate protoc -I=. -I=$GOPATH/src --gofast_out=. source.proto

func FileDescriptorBytes() []byte {
	return fileDescriptorSource
}
