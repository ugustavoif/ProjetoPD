module github.com/ruandg/microservices/shipping

go 1.26.4

replace github.com/ruandg/microservices-proto/golang/shipping => ./proto/shipping

require (
	github.com/ruandg/microservices-proto/golang/shipping v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.82.0
)

require (
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
