package resolvers

import (
	"context"
	"fmt"
	"go-backend/models"
	"strconv"

	"github.com/graph-gophers/graphql-go"
)

type MutationResolver struct {
	m models.Mutation
}

func (r *MutationResolver) CreateProduct(ctx context.Context, args struct{ Input models.ProductInput }) (*ProductResolver, error) {
	// Ensure stockQuantity is an int
	args.Input.StockQuantity = int(args.Input.StockQuantity)

	product, err := r.m.CreateProduct(args.Input)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{*product}, nil
}

func (r *MutationResolver) UpdateProduct(ctx context.Context, args struct {
	ID    graphql.ID
	Input models.ProductInput
}) (*ProductResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	product, err := r.m.UpdateProduct(id, args.Input)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{*product}, nil
}

func (r *MutationResolver) DeleteProduct(ctx context.Context, args struct{ ID graphql.ID }) (bool, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return false, fmt.Errorf("invalid ID: %v", err)
	}
	return r.m.DeleteProduct(id)
}

func (r *MutationResolver) CreateOrder(ctx context.Context, args struct{ Input models.OrderInput }) (*OrderResolver, error) {
	order, err := r.m.CreateOrder(args.Input)
	if err != nil {
		return nil, err
	}
	return &OrderResolver{*order}, nil
}

func (r *MutationResolver) CreateReview(ctx context.Context, args struct{ Input models.ReviewInput }) (*ReviewResolver, error) {
	review, err := r.m.CreateReview(args.Input)
	if err != nil {
		return nil, err
	}
	return &ReviewResolver{*review}, nil
}
