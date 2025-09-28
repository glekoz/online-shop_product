package repository

import (
	"context"
	"errors"
	"time"

	"github.com/glekoz/cache"
	"github.com/glekoz/online-shop_product/pkg/models"
	"github.com/glekoz/online-shop_product/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Если бы кэш был сложнее, мне следовало бы вынести его в отдельный пакет
// Но в моем случае кэш очень "тонкий", поэтому сделал паттерн Декоратор

type Repository struct {
	q     *db.Queries
	pool  *pgxpool.Pool
	cache *cache.Cache[string, models.Product]
}

func New(dsn string) (*Repository, error) {
	c, err := cache.New[string, models.Product]()
	if err != nil {
		return nil, err
	}
	p, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	q := db.New(p)
	return &Repository{
		q:     q,
		pool:  p,
		cache: c,
	}, nil
}

func (r *Repository) Create(ctx context.Context, id string, prod models.Product) error {
	err := r.q.Create(ctx, db.CreateParams{
		ID:          id,
		Name:        prod.Name,
		Price:       int32(prod.Price),
		Description: prod.Description},
	)
	if err != nil {
		var errp *pgconn.PgError
		if errors.As(err, &errp) {
			if errp.Code == models.UniqueErrCode {
				return models.ErrAlreadyExists
			}
		}
		return err
	}
	r.cache.Add(id, models.Product{
		Name:        prod.Name,
		Price:       prod.Price,
		Description: prod.Description,
	}, 30*time.Second)
	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (models.Product, error) {
	if res, ok := r.cache.Get(id); ok {
		return res, nil
	}
	res, err := r.q.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, models.ErrNotFound
		}
		return models.Product{}, err
	}
	return models.Product{
		Name:        res.Name,
		Price:       int(res.Price),
		Description: res.Description,
	}, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]models.ProductDigest, error) {
	ress, err := r.q.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(ress) == 0 {
		return nil, models.ErrNotFound
	}
	var result = make([]models.ProductDigest, len(ress))
	for i, res := range ress {
		result[i] = models.ProductDigest{
			ID:    res.ID,
			Name:  res.Name,
			Price: int(res.Price),
		}
	}
	return result, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	err := r.q.Delete(ctx, id)
	if err != nil {
		return err
	}
	r.cache.Delete(id)
	return nil
}

func (r *Repository) Update(ctx context.Context, id string, prod models.Product) error {
	err := r.q.Update(ctx, db.UpdateParams{
		ID:          id,
		Name:        prod.Name,
		Price:       int32(prod.Price),
		Description: prod.Description,
	})
	if err != nil {
		return err
	}
	r.cache.Add(id, prod, 30*time.Second)
	return nil
}
