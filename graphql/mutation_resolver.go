package main

import (
	"context"

	"github.com/FrankOHara43/go-grpc-microservice/order"
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	a, err := r.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	return &Account{ID: a.ID, Name: a.Name}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		return nil, err
	}
	return &Product{ID: p.ID, Name: p.Name, Description: p.Description, Price: p.Price}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	products := make([]order.OrderedProduct, 0, len(in.Products))
	for _, p := range in.Products {
		products = append(products, order.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}

	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		return nil, err
	}

	resultProducts := make([]*OrderedProduct, 0, len(o.Products))
	for _, p := range o.Products {
		resultProducts = append(resultProducts, &OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    int(p.Quantity),
		})
	}

	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
		Products:   resultProducts,
	}, nil
}
