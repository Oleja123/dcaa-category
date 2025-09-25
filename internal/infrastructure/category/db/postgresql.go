package categorydb

import (
	"context"
	"fmt"

	category "github.com/Oleja123/dcaa-category/internal/domain/category"
	"github.com/Oleja123/dcaa-property/pkg/client"
	myErrors "github.com/Oleja123/dcaa-property/pkg/errors"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	client client.Client
}

func (r *repository) Create(ctx context.Context, c category.Category) (int, error) {
	q := `
		INSERT INTO categories (category_name, info) 
		VALUES ($1, $2)
		RETURNING id
	`
	err := r.client.QueryRow(ctx, q, c.Name, c.Info).Scan(&c.Id)
	if err != nil {
		return 0, myErrors.ErrInternalError
	}
	return c.Id, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	q := `
		DELETE FROM categories WHERE id = $1
	`
	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("не найдена категория по id: %d: %w", id, myErrors.ErrNotFound)
		}
		return myErrors.ErrInternalError
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]category.Category, error) {
	q := `
		SELECT id, category_name, info FROM categories
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, myErrors.ErrInternalError
	}
	defer rows.Close()

	categories := make([]category.Category, 0)
	for rows.Next() {
		var c category.Category
		err := rows.Scan(&c.Id, &c.Name, &c.Info)
		if err != nil {
			return nil, myErrors.ErrInternalError
		}

		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, myErrors.ErrInternalError
	}

	return categories, nil
}

func (r *repository) FindOne(ctx context.Context, id int) (category.Category, error) {
	q := `
		SELECT id, category_name, info FROM categories
		WHERE id = $1
	`

	var c category.Category
	err := r.client.QueryRow(ctx, q, id).Scan(&c.Id, &c.Name, &c.Info)
	if err != nil {
		if err == pgx.ErrNoRows {
			return category.Category{}, fmt.Errorf("не найдена категория по id: %d: %w", id, myErrors.ErrNotFound)
		} else {
			return category.Category{}, myErrors.ErrInternalError
		}
	}

	return c, nil
}

func (r *repository) Update(ctx context.Context, p category.Category) error {
	q := `
		UPDATE categories SET category_name = $1, info = $2
		WHERE id = $3
	`

	_, err := r.client.Exec(ctx, q, p.Name, p.Info, p.Id)
	if err != nil {
		return myErrors.ErrInternalError
	}
	return nil
}

func NewRepository(client client.Client) category.Repository {
	return &repository{
		client: client,
	}
}
