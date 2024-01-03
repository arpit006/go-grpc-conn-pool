module grpc-test/client

go 1.21

require (
	github.com/arpit006/go-grpc-conn-pool v0.0.0-20240103204005-f734af7534cd
	//go-grpc/pool v0.0.0-20240103200111-119de504430b
	google.golang.org/grpc v1.60.1
	grpc-test/protos v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-co-op/gocron v1.37.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace grpc-test/protos => ./../protos

