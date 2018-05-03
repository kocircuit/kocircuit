package proto

//go:generate protoc -I=. --gofast_out=. types.proto

func FileDescriptorBytes() []byte {
	return fileDescriptorTypes
}
