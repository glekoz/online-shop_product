package app

import (
	"context"

	"github.com/glekoz/online-shop_product/pkg/log"
	"github.com/glekoz/online-shop_product/pkg/models"
	"github.com/google/uuid"
)

type RepoAPI interface {
	Create(ctx context.Context, id string, prod models.Product) error
	Get(ctx context.Context, id string) (models.Product, error)
	GetAll(ctx context.Context) ([]models.ProductDigest, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, prod models.Product) error
}

type App struct {
	r RepoAPI
}

func New(r RepoAPI) *App {
	return &App{r: r}
}

func (a *App) Create(ctx context.Context, prod models.Product) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		//log.WrapError()
		//return "", models.ErrInternal
		return "", err
	}
	if err = a.r.Create(ctx, uuid.String(), models.Product{
		Name:        prod.Name,
		Price:       prod.Price,
		Description: prod.Description,
	}); err != nil {
		return "", log.WrapError(ctx, err)
	}
	return uuid.String(), nil
}

func (a *App) Get(ctx context.Context, id string) (models.Product, error) {
	return a.r.Get(ctx, id)
}

func (a *App) GetAll(ctx context.Context) ([]models.ProductDigest, error) {
	return a.r.GetAll(ctx)
}

func (a *App) Delete(ctx context.Context, id string) error {
	return a.r.Delete(ctx, id)
}

// хотя нужно сделать проверку ещё в шлюзе, чтобы пустых полей в запросе (форме) не было
func (a *App) Update(ctx context.Context, id string, prod models.Product) error {
	// pr, err := a.r.Get(ctx, id)
	// if err != nil {
	// 	return err
	// }
	// if prod.Name != "" {
	// 	pr.Name = prod.Name
	// }
	// if prod.Price != 0 {
	// 	pr.Price = prod.Price
	// }
	// if prod.Description != "" {
	// 	pr.Description = prod.Description
	// }
	return a.r.Update(ctx, id, prod)
}
