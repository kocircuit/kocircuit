FROM golang:1.11-alpine

# Prepare dependencies
RUN apk add -U git

ENV GO111MODULES=auto

# Include sources
COPY . /go/src/github.com/kocircuit/kocircuit/

# Build ko
RUN \
    go get github.com/golang/protobuf/proto && \
    go get github.com/golang/protobuf/protoc-gen-go/descriptor && \
    go build -o /go/bin/ko github.com/kocircuit/kocircuit/lang/ko

# Package ko container
FROM alpine

ENV GOPATH=/ko
WORKDIR $GOPATH

# Copy binary
COPY --from=0 /go/bin/ko /usr/bin/

# Copy library sources
COPY ./lib/ /ko/src/github.com/kocircuit/kocircuit/lib/

ENTRYPOINT [ "/usr/bin/ko" ]
