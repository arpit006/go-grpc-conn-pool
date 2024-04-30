module grpc-test/server

go 1.20

require (
	golang.org/x/net v0.19.0
	google.golang.org/grpc v1.60.1
	grpc-test/protos v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace grpc-test/protos => ./../protos
