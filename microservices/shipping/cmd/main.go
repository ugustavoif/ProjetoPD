package main

import (
	"github.com/ruandg/microservices/shipping/config"
	"github.com/ruandg/microservices/shipping/internal/adapters/grpc"
	"github.com/ruandg/microservices/shipping/internal/application/core/api"
)

func main() {
	application := api.NewApplication()
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
