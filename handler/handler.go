package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/glekoz/online-shop_product/pkg/log"
	"github.com/glekoz/online-shop_product/pkg/models"
	"github.com/glekoz/online-shop_proto/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductService struct {
	app AppAPI
	product.UnimplementedGRPCProductServer
}

type AppAPI interface {
	Create(ctx context.Context, prod models.Product) (string, error)
	Get(ctx context.Context, id string) (models.Product, error)
	GetAll(ctx context.Context) ([]models.ProductDigest, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, prod models.Product) error
}

func (s *ProductService) Create(ctx context.Context, req *product.Product) (*product.ID, error) {
	prod := models.Product{
		Name:        req.GetName(),
		Price:       int(req.GetPrice()),
		Description: req.GetDescription(),
	}

	// позже эту валидацию нужно будет вынести в шлюз
	if prod.Name == "" || prod.Price <= 0 || prod.Description == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"name, price and description are required, price must be greater than 0",
		)
	}

	// все логи ниже - это прикольно, но мне ещё логировать
	// в шлюзе надо будет, из-за чего записи будут дублироваться
	ctx = log.WithProductName(ctx, prod.Name)

	// благодаря методу Handle в моем MyJSONLogHandler вся информация из контекста будет выведена в лог
	slog.InfoContext(ctx, "product creation started")

	id, err := s.app.Create(ctx, prod)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			return nil, status.Errorf(
				codes.AlreadyExists,
				"product with the same name already exists: %v", prod.Name,
			)
		}
		slog.ErrorContext(log.ErrorContext(ctx, err), "product creation: "+err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	slog.InfoContext(ctx, "product creation ended")
	return &product.ID{Id: id}, nil
}

func (s *ProductService) Get(ctx context.Context, req *product.ID) (*product.Product, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	p, err := s.app.Get(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s not found", id)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &product.Product{
		Name:        p.Name,
		Price:       int32(p.Price),
		Description: p.Description,
	}, nil
}

func (s *ProductService) GetAll(ctx context.Context, _ *emptypb.Empty) (*product.GetAllResponse, error) {
	ress, err := s.app.GetAll(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	prods := make([]*product.ProductDigest, len(ress))
	for i, res := range ress {
		prods[i] = &product.ProductDigest{
			Id:    res.ID,
			Name:  res.Name,
			Price: int32(res.Price),
		}
	}
	return &product.GetAllResponse{Products: prods}, nil
}

func (s *ProductService) Delete(ctx context.Context, req *product.ID) (*emptypb.Empty, error) {
	id := req.GetId()
	if err := s.app.Delete(ctx, id); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (s *ProductService) Update(ctx context.Context, req *product.UpdateRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	prod := req.GetProduct()

	// позже эту валидацию нужно будет вынести в шлюз
	if prod.Name == "" || prod.Price <= 0 || prod.Description == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"name, price and description are required, price must be greater than 0",
		)
	}

	if err := s.app.Update(ctx, id, models.Product{
		Name:        prod.Name,
		Price:       int(prod.Price),
		Description: prod.Description,
	}); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}
