package resolvers

import (
	"context"
	"fmt"
	"go-backend/models"
	"log"
	"strconv"

	"github.com/graph-gophers/graphql-go"
)

type QueryResolver struct {
	q models.Query
}

func (r *QueryResolver) Product(ctx context.Context, args struct{ ID graphql.ID }) (*ProductResolver, error) {
	id64, err := strconv.ParseInt(string(args.ID), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	id := int32(id64)
	product, err := r.q.Product(id)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{*product}, nil
}

func (r *QueryResolver) Products(ctx context.Context, args struct {
	Category *graphql.ID
	Search   *string
}) ([]*ProductResolver, error) {
	var categoryID *int32
	if args.Category != nil {
		id64, err := strconv.ParseInt(string(*args.Category), 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %v", err)
		}
		id := int32(id64)
		categoryID = &id
	}
	products, err := r.q.Products(categoryID, args.Search)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ProductResolver, len(products))
	for i, p := range products {
		resolvers[i] = &ProductResolver{*p}
	}
	return resolvers, nil
}

func (r *QueryResolver) Categories(ctx context.Context) ([]*CategoryResolver, error) {
	categories, err := r.q.Categories()
	if err != nil {
		return nil, err
	}
	resolvers := make([]*CategoryResolver, len(categories))
	for i, c := range categories {
		resolvers[i] = &CategoryResolver{*c}
	}
	return resolvers, nil
}

func (r *QueryResolver) Order(ctx context.Context, args struct{ ID graphql.ID }) (*OrderResolver, error) {
	log.Printf("Resolving Order for ID: %s", args.ID)
	id64, err := strconv.ParseInt(string(args.ID), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	id := int32(id64)
	order, err := r.q.Order(id)
	if err != nil {
		return nil, err
	}
	return &OrderResolver{*order}, nil
}

func (r *QueryResolver) UserOrders(ctx context.Context, args struct{ UserID graphql.ID }) ([]*OrderResolver, error) {
	userID64, err := strconv.ParseInt(string(args.UserID), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}
	userID := int32(userID64)
	orders, err := r.q.UserOrders(userID)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*OrderResolver, len(orders))
	for i, o := range orders {
		resolvers[i] = &OrderResolver{*o}
	}
	return resolvers, nil
}
