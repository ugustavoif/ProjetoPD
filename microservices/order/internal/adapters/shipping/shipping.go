package shipping_adapter

import (
	"context"
	"log"
	"time"

	"github.com/ruandg/microservices-proto/golang/shipping"
	"github.com/ruandg/microservices/order/internal/application/core/domain"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(1*time.Second)),
		)),
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) CalculateShipping(order *domain.Order) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var items []*shipping.ShippingItem
	for _, item := range order.OrderItems {
		items = append(items, &shipping.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	res, err := a.shipping.Create(ctx, &shipping.CreateShippingRequest{
		OrderId: order.ID,
		Items:   items,
	})

	if err != nil {
		log.Printf("failed to calculate shipping: %v", err)
		return 0, err
	}

	return res.DeliveryTimeDays, nil
}
