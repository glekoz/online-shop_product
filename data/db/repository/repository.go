package repository

import (
	"context"

	"github.com/glekoz/online_shop_product/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q    *Queries
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, dsn string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	queries := New(pool)
	return &Repository{q: queries, pool: pool}, nil
}

func (r *Repository) Create(ctx context.Context, p models.NewProduct, imageLimit int) (models.Product, error) { // для настройки количества изображений для каждого конкретного пользователя необходимо добавить ещё 1 параметр - кол-во фотографий, которое определялось бы в сервисе, а пока просто 10
	cp := CreateParams{
		ID:          p.ID.String(),
		UserID:      p.UserID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       int32(p.Price),
	}
	dbp, err := r.q.Create(ctx, cp)
	if err != nil {
		return models.Product{}, nil
	}
	return ExtractProductFromDB(dbp)
}

func (r *Repository) Delete(ctx context.Context, productID uuid.UUID) error {
	return r.q.Delete(ctx, productID.String())
}

func (r *Repository) DeleteByUser(ctx context.Context, userID uuid.UUID) error { // а как передавать uuid через gRPC? видимо, нужна будет строка
	return r.q.DeleteByUser(ctx, userID.String())
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Product, error) {
	res, err := r.q.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, models.ErrNotFound
	}
	prs := make([]models.Product, len(res))
	for i, r := range res {
		pr, err := ExtractProductFromDB(r)
		if err != nil {
			return nil, err
		}
		prs[i] = pr
	}
	return prs, nil
}

func (r *Repository) GetLatest(ctx context.Context) ([]models.Product, error) {
	res, err := r.q.GetLatest(ctx)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, models.ErrNotFound
	}
	prs := make([]models.Product, len(res))
	for i, r := range res {
		pr, err := ExtractProductFromDB(r)
		if err != nil {
			return nil, err
		}
		prs[i] = pr
	}
	return prs, nil
}

/*
func (r *Repository) GetOne(ctx context.Context, productID uuid.UUID) (models.FullProductInfo, error) {
	res, err := r.q.GetOne(ctx, productID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // pgx.ErrNoRows wraps sql.ErrNoRows
			return models.FullProductInfo{}, models.ErrNotFound
		}
		return models.FullProductInfo{}, err
	}
	pr, err := ExtractFullProductInfoFromDB(res)
	if err != nil {
		return models.FullProductInfo{}, err
	}
	return pr, nil
}
*/
