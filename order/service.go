package order

import (
	"context"
	"time"
)

type service interface {
	PostOrder
	GetOrdersForAccount
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type OrderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {

}

func (s orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {

}
