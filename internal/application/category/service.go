package categoryservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Oleja123/dcaa-category/internal/domain/category"
	categoryDto "github.com/Oleja123/dcaa-property/pkg/dto/category"
	myErrors "github.com/Oleja123/dcaa-property/pkg/errors"
	optionalType "github.com/denpa16/optional-go-type"
)

type categoryService struct {
	repository category.Repository
}

func (cs *categoryService) CategoryToDTO(p category.Category) categoryDto.CategoryDTO {
	dto := categoryDto.CategoryDTO{}
	dto.Name = optionalType.NewOptionalString(&p.Name)
	if !p.Info.Valid {
		dto.Info = optionalType.NewOptionalString(nil)
	} else {
		dto.Info = optionalType.NewOptionalString(&p.Info.String)
	}
	dto.Id = optionalType.NewOptionalInt(&p.Id)
	return dto
}

func (cs *categoryService) CategoryFromDTO(ctx context.Context, dto categoryDto.CategoryDTO) category.Category {
	p := category.Category{}
	p.Name = *dto.Name.Value
	if dto.Id.Valid {
		p.Id = *dto.Id.Value
	}
	if dto.Info.Valid {
		p.Info = sql.NullString{
			String: *dto.Info.Value,
			Valid:  true,
		}
	}
	return p
}

func (cs *categoryService) Create(ctx context.Context, dto categoryDto.CategoryDTO) (int, error) {
	category := cs.CategoryFromDTO(ctx, dto)
	id, err := cs.repository.Create(ctx, category)
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("ошибка при создании сущности категории: %w", err)
	}
	return id, nil
}

func (cs *categoryService) Update(ctx context.Context, dto categoryDto.CategoryDTO) error {
	if _, err := cs.repository.FindOne(ctx, *dto.Id.Value); err != nil {
		return fmt.Errorf("ошибка при обновлении сущности с id: %d: %w", *dto.Id.Value, err)
	}
	pr := cs.CategoryFromDTO(ctx, dto)
	err := cs.repository.Update(ctx, pr)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ошибка при обновлении сущности категории с id: %d: %w", *dto.Id.Value, err)
	}
	return nil
}

func (cs *categoryService) Delete(ctx context.Context, id int) error {
	if _, err := cs.repository.FindOne(ctx, id); err != nil {
		return fmt.Errorf("ошибка при удалении сущности с id: %d: %w", id, err)
	}
	err := cs.repository.Delete(ctx, id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ошибка при удалении сущности категории с id: %d: %w", id, err)
	}
	return nil
}

func (cs *categoryService) FindAll(ctx context.Context) ([]categoryDto.CategoryDTO, error) {
	ca, err := cs.repository.FindAll(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("ошибка при получении списка записей категорий: %w", err)
	}
	res := make([]categoryDto.CategoryDTO, 0, len(ca))
	for _, val := range ca {
		res = append(res, cs.CategoryToDTO(val))
	}

	return res, nil
}

func (cs *categoryService) FindOne(ctx context.Context, id int) (categoryDto.CategoryDTO, error) {
	ca, err := cs.repository.FindOne(ctx, id)
	switch {
	case errors.Is(err, myErrors.ErrNotFound):
		return categoryDto.CategoryDTO{}, fmt.Errorf("ошибка при получении записи категории с id: %d: %w", id, err)
	case err != nil:
		return categoryDto.CategoryDTO{}, fmt.Errorf("ошибка при получении записи категории с id: %d: %w", id, err)
	}

	return cs.CategoryToDTO(ca), nil
}

func NewService(repo category.Repository) category.Service {
	return &categoryService{
		repository: repo,
	}
}
