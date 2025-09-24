package category

import (
	"context"

	otherService "github.com/Oleja123/dcaa-property/pkg/dto/category"
)

type Service interface {
	Create(ctx context.Context, dto otherService.CategoryDTO) (int, error)
	FindAll(ctx context.Context) ([]otherService.CategoryDTO, error)
	FindOne(ctx context.Context, id int) (otherService.CategoryDTO, error)
	Update(ctx context.Context, category otherService.CategoryDTO) error
	Delete(ctx context.Context, id int) error
}
