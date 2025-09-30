package handler

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/glekoz/online-shop_product/pkg/models"
	"github.com/glekoz/online-shop_proto/product"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// ----------------------------------------------------------------
// 						TEST SUITE SETUP SECTION
// ----------------------------------------------------------------

type ServerSuite struct {
	suite.Suite
	conn   *grpc.ClientConn
	client product.GRPCProductClient
	ctx    context.Context
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func (s *ServerSuite) SetupSuite() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		go NewServer(&AppMock{}).RunServer(8000)
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}()
	wg.Wait()
	ctx := context.Background()
	conn, err := grpc.NewClient("127.0.0.1:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := product.NewGRPCProductClient(conn)
	s.ctx = ctx
	s.client = client
	s.conn = conn
}

func (s *ServerSuite) TearDownSuite() {
	s.conn.Close()
}

// ----------------------------------------------------------------
// 							MOCK SECTION
// ----------------------------------------------------------------

type AppMock struct {
}

func (a *AppMock) Create(ctx context.Context, prod models.Product) (string, error) {
	if prod.Name == "Donut" {
		return "", models.ErrAlreadyExists
	} else if prod.Name == "Unknown" {
		return "", models.ErrInternal
	}
	return "10", nil
}

func (a *AppMock) Get(ctx context.Context, id string) (models.Product, error) {
	if id == "500" {
		return models.Product{}, models.ErrInternal
	} else if id == "404" {
		return models.Product{}, models.ErrNotFound
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
	if id == "500" {
		return models.ErrInternal
	} else if id == "404" {
		return models.ErrNotFound
	}
	return nil
}

func (a *AppMock) Update(ctx context.Context, id string, prod models.Product) error {
	if id == "500" {
		return models.ErrInternal
	} else if id == "404" {
		return models.ErrNotFound
	}
	return nil
}

func (s *ServerSuite) TestCreate() {
	tests := []struct {
		name       string
		prod       *product.Product
		expectedID string
		errCode    codes.Code
		errMsg     string
	}{
		{
			name:       "Happy",
			prod:       &product.Product{Name: "Tasty Donut", Price: 1000, Description: "Tasty"},
			expectedID: "10",
			errCode:    codes.OK,
			errMsg:     "",
		},
		{
			name:       "Already Exists",
			prod:       &product.Product{Name: "Donut", Price: 1111, Description: "Tasty"},
			expectedID: "",
			errCode:    codes.AlreadyExists,
			errMsg:     "product with the same name already exists: Donut",
		},
		{
			name:       "Internal Error",
			prod:       &product.Product{Name: "Unknown", Price: 88888888, Description: "???"},
			expectedID: "",
			errCode:    codes.Internal,
			errMsg:     models.ErrInternal.Error(),
		},
		{
			name:       "Invalid Argument",
			prod:       &product.Product{},
			expectedID: "",
			errCode:    codes.InvalidArgument,
			errMsg:     "name, price and description are required, price must be greater than 0",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			response, err := s.client.Create(s.ctx, tt.prod)
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

func (s *ServerSuite) TestGet() {
	tests := []struct {
		name         string
		id           string
		expectedProd *product.Product
		errCode      codes.Code
		errMsg       string
	}{
		{
			name:         "Happy",
			id:           "1",
			expectedProd: &product.Product{Name: "Donut", Price: 1000, Description: "Delicious"},
			errCode:      codes.OK,
			errMsg:       "",
		},
		{
			name:         "Invalid Argument",
			id:           "",
			expectedProd: &product.Product{},
			errCode:      codes.InvalidArgument,
			errMsg:       "id is required",
		},
		{
			name:         "Not Found",
			id:           "404",
			expectedProd: &product.Product{},
			errCode:      codes.NotFound,
			errMsg:       "404 not found",
		},
		{
			name:         "Internal",
			id:           "500",
			expectedProd: &product.Product{},
			errCode:      codes.Internal,
			errMsg:       models.ErrInternal.Error(),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			response, err := s.client.Get(s.ctx, &product.ID{Id: tt.id})
			if response != nil {
				if response.GetName() != tt.expectedProd.Name {
					t.Error("response: expected: ", tt.expectedProd.Name, "received: ", response.GetName())
				}
				if response.GetPrice() != tt.expectedProd.Price {
					t.Error("response: expected: ", tt.expectedProd.Price, "received: ", response.GetPrice())
				}
				if response.GetDescription() != tt.expectedProd.Description {
					t.Error("response: expected: ", tt.expectedProd.Description, "received: ", response.GetDescription())
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
