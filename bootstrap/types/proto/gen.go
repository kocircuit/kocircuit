package proto

//go:generate protoc -I=. -I=$GOPATH/src --gofast_out=. types.proto

func FileDescriptorBytes() []byte {
	return fileDescriptorTypes
}
