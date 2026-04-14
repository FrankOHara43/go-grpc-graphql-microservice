package main

import "context"

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	orders, err := r.server.orderClient.GetOrdersForAccount(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	result := make([]*Order, 0, len(orders))
	for _, o := range orders {
		products := make([]*OrderedProduct, 0, len(o.Products))
		for _, p := range o.Products {
			products = append(products, &OrderedProduct{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
			})
		}

		result = append(result, &Order{
			ID:         o.ID,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
			Products:   products,
		})
	}

	return result, nil
}
