package categoryservice_test

import (
	"context"
	"testing"

	categoryservice "github.com/Oleja123/dcaa-category/internal/application/category"
	categorymock "github.com/Oleja123/dcaa-category/internal/infrastructure/category"
	categorydto "github.com/Oleja123/dcaa-property/pkg/dto/category"
	myErrors "github.com/Oleja123/dcaa-property/pkg/errors"
	optionalType "github.com/denpa16/optional-go-type"
	"github.com/stretchr/testify/assert"
)

func TestCategoryService(t *testing.T) {
	ctx := context.Background()
	repo := &categorymock.MockCategoryRepo{}
	service := categoryservice.NewService(repo)

	name := "Category"
	info := "Cool"

	dto := categorydto.CategoryDTO{
		Name: optionalType.NewOptionalString(&name),
		Info: optionalType.NewOptionalString(&info),
	}
	id, err := service.Create(ctx, dto)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	dto.Id = optionalType.NewOptionalInt(&id)
	err = service.Update(ctx, dto)
	assert.NoError(t, err)

	id = 0
	dto.Id = optionalType.NewOptionalInt(&id)
	err = service.Update(ctx, dto)
	assert.ErrorAs(t, err, &myErrors.ErrNotFound)

	err = service.Delete(ctx, 1)
	assert.NoError(t, err)

	err = service.Delete(ctx, 0)
	assert.ErrorAs(t, err, &myErrors.ErrNotFound)

	all, err := service.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, all, 1)
	assert.Equal(t, "Category", *all[0].Name.Value)

	found, err := service.FindOne(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Category", *found.Name.Value)

	_, err = service.FindOne(ctx, 0)
	assert.ErrorAs(t, err, &myErrors.ErrNotFound)
}
