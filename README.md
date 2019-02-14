# protobuf-sample

sample of protocol buffer with golang

this sample implements two api by using protobuf:
1. adduser
2. getuser  


the protobuf is used to implement the api in 3ways:

1. use protobuf as type definition tool (normal.go)
2. use protobuf for rpc (rpc.go)
3. use protobuf for grpc (gprc.go)
 
## install protoc

1. install brotobuf-all-x.x.x from [protobuf](https://developers.google.com/protocol-buffers/docs/downloads)
2. unzip zip
3. cd protobuf-x.x.x
4. ./configure
5. make
6. make check
6. sudo make install

## install golang protocol buffer plugin

`go get -u github.com/golang/protobuf/protoc-gen-go`

## compile .proto

`protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto`

## for more info

tutorial: https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers
doc: https://developers.google.com/protocol-buffers/docs/proto3