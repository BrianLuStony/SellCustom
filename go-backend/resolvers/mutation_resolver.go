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
	categoryID, err := strconv.ParseInt(string(args.Input.CategoryID), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %v", err)
	}

	// Create a new ProductInput with the correct types
	input := models.ProductInput{
		Name:          args.Input.Name,
		Description:   args.Input.Description,
		Price:         args.Input.Price,
		StockQuantity: args.Input.StockQuantity,
		CategoryID:    int32(categoryID),
		Images:        args.Input.Images,
	}

	product, err := r.m.CreateProduct(input)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{*product}, nil
}

func (r *MutationResolver) UpdateProduct(ctx context.Context, args struct {
	ID    graphql.ID
	Input models.ProductInput
}) (*ProductResolver, error) {
	id64, err := strconv.ParseInt(string(args.ID), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	id := int32(id64)
	product, err := r.m.UpdateProduct(id, args.Input)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{*product}, nil
}

func (r *MutationResolver) DeleteProduct(ctx context.Context, args struct{ ID graphql.ID }) (bool, error) {
	id64, err := strconv.ParseInt(string(args.ID), 10, 32)
	if err != nil {
		return false, fmt.Errorf("invalid ID: %v", err)
	}
	id := int32(id64)
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

func (r *MutationResolver) CreateCategory(ctx context.Context, args struct{ Input models.CategoryInput }) (*CategoryResolver, error) {
	var parentID *int32
	if args.Input.ParentID != nil {
		id, err := strconv.ParseInt(*args.Input.ParentID, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid parent ID: %v", err)
		}
		parentIDInt32 := int32(id)
		parentID = &parentIDInt32
	}

	// Use the new method name here
	category, err := r.m.CreateNewCategory(args.Input.Name, parentID)
	if err != nil {
		return nil, err
	}
	return &CategoryResolver{*category}, nil
}
