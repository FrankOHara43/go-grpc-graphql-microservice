package main

import "context"

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	if id != nil {
		a, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Account{{ID: a.ID, Name: a.Name}}, nil
	}

	skip := uint64(0)
	take := uint64(20)
	if pagination != nil {
		if pagination.Skip != nil && *pagination.Skip >= 0 {
			skip = uint64(*pagination.Skip)
		}
		if pagination.Take != nil && *pagination.Take >= 0 {
			take = uint64(*pagination.Take)
		}
	}

	accounts, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	result := make([]*Account, 0, len(accounts))
	for _, a := range accounts {
		result = append(result, &Account{ID: a.ID, Name: a.Name})
	}
	return result, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	if id != nil {
		p, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Product{{ID: p.ID, Name: p.Name, Description: p.Description, Price: p.Price}}, nil
	}

	skip := uint64(0)
	take := uint64(20)
	if pagination != nil {
		if pagination.Skip != nil && *pagination.Skip >= 0 {
			skip = uint64(*pagination.Skip)
		}
		if pagination.Take != nil && *pagination.Take >= 0 {
			take = uint64(*pagination.Take)
		}
	}

	q := ""
	if query != nil {
		q = *query
	}

	products, err := r.server.catalogClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		return nil, err
	}

	result := make([]*Product, 0, len(products))
	for _, p := range products {
		result = append(result, &Product{ID: p.ID, Name: p.Name, Description: p.Description, Price: p.Price})
	}
	return result, nil
}
