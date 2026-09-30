package domain

import (
	"time"
)

type ShippingItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int32  `json:"quantity"`
}

type Shipping struct {
	ID           int64          `json:"id"`
	CustomerID   int64          `json:"customer_id"`
	OrderID      int64          `json:"order_id"`
	Items        []ShippingItem `json:"items"`
	DeliveryDays int32          `json:"delivery_days"`
	CreatedAt    int64          `json:"created_at"`
}

func NewShipping(customerId int64, orderId int64, items []ShippingItem) Shipping {
	return Shipping{
		CreatedAt:  time.Now().Unix(),
		CustomerID: customerId,
		OrderID:    orderId,
		Items:      items,
	}
}

// EstimateDeliveryDays is one day, plus one more for every 5 items.
func (s Shipping) EstimateDeliveryDays() int32 {
	var count int32
	for _, item := range s.Items {
		count += item.Quantity
	}
	return 1 + count/5
}
