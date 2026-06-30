module github.com/ruandg/microservices/order

go 1.26.4

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/ruandg/microservices-proto/golang/order v0.0.0-00010101000000-000000000000
	github.com/ruandg/microservices-proto/golang/payment v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.82.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/ruandg/microservices-proto/golang/shipping v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/ruandg/microservices-proto/golang/order => ./proto/order

replace github.com/ruandg/microservices-proto/golang/payment => ./proto/payment

replace github.com/ruandg/microservices-proto/golang/shipping => ./proto/shipping
