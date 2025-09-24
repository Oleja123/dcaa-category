package category

import "context"

type Repository interface {
	Create(ctx context.Context, category Category) (int, error)
	FindAll(ctx context.Context) ([]Category, error)
	FindOne(ctx context.Context, id int) (Category, error)
	Update(ctx context.Context, category Category) error
	Delete(ctx context.Context, id int) error
}
