package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstimateDeliveryDays(t *testing.T) {
	tests := []struct {
		name  string
		items []ShippingItem
		want  int32
	}{
		{"no items", nil, 1},
		{"under five", []ShippingItem{{ProductCode: "CAM", Quantity: 4}}, 1},
		{"exactly five", []ShippingItem{{ProductCode: "CAM", Quantity: 5}}, 2},
		{"summed across items", []ShippingItem{{ProductCode: "CAM", Quantity: 3}, {ProductCode: "BAG", Quantity: 8}}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShipping(1, 1, tt.items)
			assert.Equal(t, tt.want, s.EstimateDeliveryDays())
		})
	}
}
