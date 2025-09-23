package repository

import (
	"context"
	"errors"

	"github.com/glekoz/online-shop_product/models"
	"github.com/glekoz/online-shop_product/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q    *db.Queries
	pool *pgxpool.Pool
}

func New(dsn string) (*Repository, error) {
	p, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	q := db.New(p)
	return &Repository{
		q:    q,
		pool: p,
	}, nil
}

func (r *Repository) Create(ctx context.Context, args db.CreateParams) error {
	return r.q.Create(ctx, args)
}

func (r *Repository) Get(ctx context.Context, id string) (db.GetRow, error) {
	res, err := r.q.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.GetRow{}, models.ErrNotFound
		}
		return db.GetRow{}, err
	}
	return res, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]db.GetAllRow, error) {
	ress, err := r.q.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(ress) == 0 {
		return nil, models.ErrNotFound
	}
	return ress, nil
}
