package domain

type ShippingItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int32  `json:"quantity"`
}

type Shipping struct {
	OrderID          int64          `json:"order_id"`
	Items            []ShippingItem `json:"items"`
	DeliveryTimeDays int32          `json:"delivery_time_days"`
}

func NewShipping(orderId int64, items []ShippingItem) Shipping {
	var totalItems int32
	for _, item := range items {
		totalItems += item.Quantity
	}

	deliveryTime := int32(1) + (totalItems / 5)

	return Shipping{
		OrderID:          orderId,
		Items:            items,
		DeliveryTimeDays: deliveryTime,
	}
}
