package handler

import (
	"fmt"
	"net"

	"github.com/glekoz/online-shop_proto/product"
	"google.golang.org/grpc"
)

func NewServer(app AppAPI) *ProductService {
	return &ProductService{app: app}
}

func (ps *ProductService) RunServer(port int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	product.RegisterGRPCProductServer(serv, ps)
	return serv.Serve(listen)
}
