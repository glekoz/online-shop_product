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

// ----------------------------------------------------------------
// 							TEST SECTION
// ----------------------------------------------------------------

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
				s.Assert().Equal(response.GetId(), tt.expectedID)
			}
			if err != nil {
				if er, ok := status.FromError(err); ok {
					s.Assert().Equal(er.Code(), tt.errCode)
					s.Assert().Equal(er.Message(), tt.errMsg)
				}
			}
		})
	}
}

func (s *ServerSuite) TestGet() {
	tests := []struct {
		name                string
		id                  string
		expectedName        string
		expectedPrice       int32
		expectedDescription string
		errCode             codes.Code
		errMsg              string
	}{
		{
			name:                "Happy",
			id:                  "1",
			expectedName:        "Donut",
			expectedPrice:       1000,
			expectedDescription: "Delicious",
			errCode:             codes.OK,
			errMsg:              "",
		},
		{
			name:                "Invalid Argument",
			id:                  "",
			expectedName:        "",
			expectedPrice:       0,
			expectedDescription: "",
			errCode:             codes.InvalidArgument,
			errMsg:              "id is required",
		},
		{
			name:                "Not Found",
			id:                  "404",
			expectedName:        "",
			expectedPrice:       0,
			expectedDescription: "",
			errCode:             codes.NotFound,
			errMsg:              "404 not found",
		},
		{
			name:                "Internal",
			id:                  "500",
			expectedName:        "",
			expectedPrice:       0,
			expectedDescription: "",
			errCode:             codes.Internal,
			errMsg:              models.ErrInternal.Error(),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			response, err := s.client.Get(s.ctx, &product.ID{Id: tt.id})
			if response != nil {
				s.Assert().Equal(response.GetName(), tt.expectedName)
				s.Assert().Equal(response.GetPrice(), tt.expectedPrice)
				s.Assert().Equal(response.GetDescription(), tt.expectedDescription)
			}
			if err != nil {
				if er, ok := status.FromError(err); ok {
					s.Assert().Equal(er.Code(), tt.errCode)
					s.Assert().Equal(er.Message(), tt.errMsg)
				}
			}
		})
	}
}

func (s *ServerSuite) TestGetAll() {
	tests := []struct {
		name           string
		expectedIDs    []string
		expectedNames  []string
		expectedPrices []int32
		errCode        codes.Code
		errMsg         string
	}{
		{
			name:           "Happy",
			expectedIDs:    []string{"1", "2", "3"},
			expectedNames:  []string{"Donut", "Another Donut", "Another Another Donut"},
			expectedPrices: []int32{1000, 1200, 1500},
			errCode:        codes.OK,
			errMsg:         "",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			response, err := s.client.GetAll(s.ctx, nil)
			if response != nil {
				for i, res := range response.GetProducts() {
					s.Assert().Equal(tt.expectedIDs[i], res.GetId())
					s.Assert().Equal(tt.expectedNames[i], res.GetName())
					s.Assert().Equal(tt.expectedPrices[i], res.GetPrice())
				}
			}
			s.Assert().Nil(err)
		})
	}
}

func (s *ServerSuite) TestDelete() {
	// if id == "500" {
	// 	return models.ErrInternal
	// } else if id == "404" {
	// 	return models.ErrNotFound
	// }
	tests := []struct {
		name    string
		id      string
		errCode codes.Code
		errMsg  string
	}{
		{
			name:    "Happy",
			id:      "1",
			errCode: codes.OK,
			errMsg:  "",
		},
		{
			name:    "Not Found",
			id:      "404",
			errCode: codes.NotFound,
			errMsg:  models.ErrNotFound.Error(),
		},
		{
			name:    "Internal",
			id:      "500",
			errCode: codes.Internal,
			errMsg:  models.ErrInternal.Error(),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			resp, err := s.client.Delete(s.ctx, &product.ID{Id: tt.id})
			if err != nil {
				s.Assert().Nil(resp) // gRPC возвращает nil, если произошла ошибка
				er, _ := status.FromError(err)
				s.Assert().Equal(er.Code(), tt.errCode)
				s.Assert().Equal(er.Message(), tt.errMsg)
			} else {
				s.Assert().NotNil(resp) // gRPC возвращает не nil, если ошибки не было
			}
		})
	}
}

func (s *ServerSuite) TestUpdate() {
	tests := []struct {
		name    string
		id      string
		prod    models.Product
		errCode codes.Code
		errMsg  string
	}{
		{
			name:    "Happy",
			id:      "1",
			prod:    models.Product{Name: "Donut", Price: 1000, Description: "Tasty"},
			errCode: codes.OK,
			errMsg:  "",
		},
		{
			name:    "Invalid Argument",
			id:      "2",
			prod:    models.Product{},
			errCode: codes.InvalidArgument,
			errMsg:  "name, price and description are required, price must be greater than 0",
		},
		{
			name:    "Not Found",
			id:      "404",
			prod:    models.Product{Name: "Donut", Price: 1000, Description: "Tasty"},
			errCode: codes.NotFound,
			errMsg:  models.ErrNotFound.Error(),
		},
		{
			name:    "Internal",
			id:      "500",
			prod:    models.Product{Name: "Donut", Price: 1000, Description: "Tasty"},
			errCode: codes.Internal,
			errMsg:  models.ErrInternal.Error(),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			_, err := s.client.Update(s.ctx, &product.UpdateRequest{
				Id: tt.id,
				Product: &product.Product{
					Name:        tt.prod.Name,
					Price:       int32(tt.prod.Price),
					Description: tt.prod.Description,
				}})
			er, _ := status.FromError(err)
			s.Assert().Equal(tt.errCode, er.Code())
			s.Assert().Equal(tt.errMsg, er.Message())
		})
	}
}
