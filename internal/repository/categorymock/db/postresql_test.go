package categorydb_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Oleja123/dcaa-category/internal/domain/category"
	categorydb "github.com/Oleja123/dcaa-category/internal/repository/categorymock/db"
	"github.com/Oleja123/dcaa-property/pkg/client"
	"github.com/Oleja123/dcaa-property/pkg/client/postgresql"
	"github.com/Oleja123/dcaa-property/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) client.Client {
	cfg, err := config.LoadConfig("../../../../test_config .yaml")
	require.NoError(t, err)

	client, err := postgresql.NewClient(t.Context(), cfg)
	require.NoError(t, err)

	_, err = client.Exec(context.Background(), "TRUNCATE TABLE properties RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	return client
}

func TestPropertyRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	client := setupTestDB(t)
	repo := categorydb.NewRepository(client)

	p := category.Category{
		Name: "Category 1",
		Info: sql.NullString{String: "Nice Category", Valid: true},
	}

	id, err := repo.Create(ctx, p)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.FindOne(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Category 1", got.Name)
	assert.Equal(t, "Nice Category", got.Info.String)

	got.Name = "Updated Category"
	err = repo.Update(ctx, got)
	require.NoError(t, err)

	updated, err := repo.FindOne(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Updated Category", updated.Name)

	all, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Len(t, all, 1)

	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	_, err = repo.FindOne(ctx, id)
	assert.Error(t, err)
}
