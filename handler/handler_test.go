package handler

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/glekoz/online-shop_product/pkg/models"
	"github.com/glekoz/online-shop_proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type AppMock struct {
}

func (a *AppMock) Create(ctx context.Context, prod models.Product) (string, error) {
	if prod.Name == "" || prod.Price < 0 || prod.Description == "" {
		return "", status.Error(codes.InvalidArgument, "price must be > 0, name and description are required")
	}
	return "1", nil
}

func (a *AppMock) Get(ctx context.Context, id string) (models.Product, error) {
	if id == "0" {
		return models.Product{}, status.Error(codes.NotFound, "there is no entry with id = 0")
	}
	return models.Product{Name: "Donut", Price: 1000, Description: "Delicious"}, nil
}

func (a *AppMock) GetAll(ctx context.Context) ([]models.ProductDigest, error) {
	return []models.ProductDigest{
		{ID: "1", Name: "Donut", Price: 1000},
		{ID: "2", Name: "Another Donut", Price: 1200},
		{ID: "3", Name: "Another Another Donut", Price: 1500},
	}, nil
}

func (a *AppMock) Delete(ctx context.Context, id string) error {
	if id == "0" {
		return status.Error(codes.NotFound, "there is no entry with id = 0")
	}
	return nil
}

func (a *AppMock) Update(ctx context.Context, id string, prod models.Product) error {
	if id == "0" {
		return status.Error(codes.NotFound, "there is no entry with id = 0")
	}
	if prod.Name == "" || prod.Price == 0 || prod.Description == "" {
		return status.Error(codes.InvalidArgument, "price must be > 0, name and description are required")
	}
	return nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	product.RegisterGRPCProductServer(server, NewServer(&AppMock{}))

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestProductServer(t *testing.T) {
	tests := []struct {
		name       string
		prod       models.Product
		expectedID string
		errCode    codes.Code
		errMsg     string
	}{
		{
			name:       "Happy",
			prod:       models.Product{Name: "Donut", Price: 1000, Description: "Tasty"},
			expectedID: "1",
			errCode:    codes.OK,
			errMsg:     "",
		},
		{
			name:       "Error",
			prod:       models.Product{},
			expectedID: "",
			errCode:    codes.InvalidArgument,
			errMsg:     "price must be > 0, name and description are required",
		},
	}

	ctx := context.Background()
	conn, err := grpc.NewClient("", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := product.NewGRPCProductClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &product.Product{Name: tt.prod.Name, Price: int32(tt.prod.Price), Description: tt.prod.Description}
			response, err := client.Create(ctx, request)
			if response != nil {
				if response.GetId() != tt.expectedID {
					t.Error("response: expected: ", tt.expectedID, "received: ", response.GetId())
				}
			}
			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected: ", tt.errCode, "received: ", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected: ", tt.errMsg, "received: ", er.Message())
					}
				}
			}
		})
	}
}
