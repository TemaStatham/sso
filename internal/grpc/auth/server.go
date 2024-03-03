package auth

import (
	"context"

	ssov1 "github.com/TemaStatham/protos/gen/go/sso"
	"google.golang.org/grpc"
)

// serverAPI реализует  интерфейсы обработки запросов для работы с grpc
type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

// Register регистрирует обработчик в grpc
func Register(grpc *grpc.Server) {
	ssov1.RegisterAuthServer(grpc, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{Token: "1234"}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterResquest) (*ssov1.RegisterResponse, error) {
	panic("impement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("impement me")
}
