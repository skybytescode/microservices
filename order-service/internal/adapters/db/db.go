package db

import (
	"fmt"

	"github.com/skybytescode/microservices/order-service/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID     uint
	ProductCode string
	UnitPrice   float32
	Quantity    int32
}

type Adapter struct{ db *gorm.DB }

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, err := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}
	if err := db.AutoMigrate(&Order{}, &OrderItem{}); err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a Adapter) Save(order *domain.Order) error {
	orderItems := make([]OrderItem, 0, len(order.OrderItems))
	for _, it := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: it.ProductCode,
			UnitPrice:   it.UnitPrice,
			Quantity:    it.Quantity,
		})
	}
	orderModel := Order{CustomerID: order.CustomerID, Status: order.Status, OrderItems: orderItems}
	res := a.db.Create(&orderModel)
	if res.Error == nil {
		order.ID = int64(orderModel.ID)
	}
	return res.Error
}

func (a Adapter) Get(id string) (domain.Order, error) {
	var orderModel Order
	res := a.db.Preload("OrderItems").First(&orderModel, "id = ?", id)
	if res.Error != nil {
		return domain.Order{}, res.Error
	}

	items := make([]domain.OrderItem, 0, len(orderModel.OrderItems))
	for _, it := range orderModel.OrderItems {
		items = append(items, domain.OrderItem{
			ProductCode: it.ProductCode,
			UnitPrice:   it.UnitPrice,
			Quantity:    it.Quantity,
		})
	}

	return domain.Order{
		ID:         int64(orderModel.ID),
		CustomerID: orderModel.CustomerID,
		Status:     orderModel.Status,
		OrderItems: items,
	}, nil
}
