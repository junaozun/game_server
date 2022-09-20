#nats的protobuf使用的google的protocbuf，所以这里需要安装google的protoc-gen-go，而不能用github的protoc-gen-go插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26

protoc --proto_path=./ --go_out=./ my.proto

