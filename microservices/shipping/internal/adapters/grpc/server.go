package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ruandg/microservices-proto/golang/shipping"
	"github.com/ruandg/microservices/shipping/config"
	"github.com/ruandg/microservices/shipping/internal/application/core/domain"
	"github.com/ruandg/microservices/shipping/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api  ports.APIPort
	port int
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Create(ctx context.Context, request *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	log.Printf("====================================================")
	log.Printf("Calculando prazo de envio para o pedido %d", request.OrderId)
	var shippingItems []domain.ShippingItem
	for _, item := range request.Items {
		shippingItems = append(shippingItems, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	newShipping := domain.NewShipping(request.OrderId, shippingItems)

	result, err := a.api.CalculateShipping(newShipping)
	if err != nil {
		return nil, err
	}

	log.Printf("Prazo de envio calculado para o pedido %d: %d dias", request.OrderId, result.DeliveryTimeDays)
	log.Printf("====================================================")

	return &shipping.CreateShippingResponse{DeliveryTimeDays: result.DeliveryTimeDays}, nil
}

func (a Adapter) Run() {
	var err error
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("Falha ao escutar na porta %d, erro: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	shipping.RegisterShippingServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Falha ao iniciar servidor grpc na porta %d", a.port)
	}
}
