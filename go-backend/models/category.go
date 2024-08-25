// models/category.go

package models

import (
	"go-backend/db"
)

func (m *Mutation) CreateNewCategory(name string, parentID *int32) (*Category, error) {
	var category Category
	var err error

	if parentID != nil {
		err = db.DB.QueryRow(`
            INSERT INTO categories (name, parent_id)
            VALUES ($1, $2)
            RETURNING id, name, parent_id
        `, name, parentID).Scan(&category.ID, &category.Name, &category.ParentCategory.ID)
	} else {
		err = db.DB.QueryRow(`
            INSERT INTO categories (name)
            VALUES ($1)
            RETURNING id, name
        `, name).Scan(&category.ID, &category.Name)
	}

	if err != nil {
		return nil, err
	}

	if parentID != nil {
		category.ParentCategory = &Category{ID: *parentID}
	}

	return &category, nil
}
