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
		prod       models.Product
		expectedID string
		errCode    codes.Code
		errMsg     string
	}{
		{
			name:       "Happy",
			prod:       models.Product{Name: "Tasty Donut", Price: 1000, Description: "Tasty"},
			expectedID: "10",
			errCode:    codes.OK,
			errMsg:     "",
		},
		{
			name:       "Already Exists",
			prod:       models.Product{Name: "Donut", Price: 1111, Description: "Tasty"},
			expectedID: "",
			errCode:    codes.AlreadyExists,
			errMsg:     "product with the same name already exists: Donut",
		},
		{
			name:       "Internal Error",
			prod:       models.Product{Name: "Unknown", Price: 88888888, Description: "???"},
			expectedID: "",
			errCode:    codes.Internal,
			errMsg:     models.ErrInternal.Error(),
		},
		{
			name:       "Invalid Argument",
			prod:       models.Product{},
			expectedID: "",
			errCode:    codes.InvalidArgument,
			errMsg:     "name, price and description are required, price must be greater than 0",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			request := &product.Product{Name: tt.prod.Name, Price: int32(tt.prod.Price), Description: tt.prod.Description}
			response, err := s.client.Create(s.ctx, request)
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
