package categorymock

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Oleja123/dcaa-category/internal/domain/category"
)

type MockCategoryRepo struct{}

func (m *MockCategoryRepo) Create(ctx context.Context, p category.Category) (int, error) {
	if p.Name == "fail" {
		return 0, errors.New("creation failed")
	}
	return 1, nil
}

func (m *MockCategoryRepo) Update(ctx context.Context, p category.Category) error {
	if p.Id == 0 {
		return errors.New("update failed")
	}
	return nil
}

func (m *MockCategoryRepo) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("delete failed")
	}
	return nil
}

func (m *MockCategoryRepo) FindAll(ctx context.Context) ([]category.Category, error) {
	return []category.Category{
		{
			Id:   1,
			Name: "Category",
			Info: sql.NullString{String: "Nice category", Valid: true},
		},
	}, nil
}

func (m *MockCategoryRepo) FindOne(ctx context.Context, id int) (category.Category, error) {
	if id == 0 {
		return category.Category{}, errors.New("not found")
	}
	return category.Category{
		Id:   id,
		Name: "Category",
		Info: sql.NullString{String: "Nice category", Valid: true},
	}, nil
}
