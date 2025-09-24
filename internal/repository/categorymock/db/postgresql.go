package categorydb

import (
	"context"
	"fmt"

	category "github.com/Oleja123/dcaa-category/internal/domain/category"
	"github.com/Oleja123/dcaa-property/pkg/client"
	"github.com/jackc/pgx/v5/pgconn"
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
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return 0, fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
		}
		return 0, err
	}
	return c.Id, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	q := `
		DELETE FROM categories WHERE id = $1
	`
	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
		}
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]category.Category, error) {
	q := `
		SELECT id, category_name, info FROM categories
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]category.Category, 0)
	for rows.Next() {
		var c category.Category
		err := rows.Scan(&c.Id, &c.Name, &c.Info)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		return category.Category{}, err
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
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
		}
		return err
	}
	return nil
}

func NewRepository(client client.Client) category.Repository {
	return &repository{
		client: client,
	}
}
