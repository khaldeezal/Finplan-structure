package grpcclients

import (
	"context"
	"os"
	"time"

	"github.com/khaldeezal/Finplan-proto/proto-definitions/gen/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient базовый gRPC-клиент для AuthService.
type AuthClient struct {
	client auth.AuthServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// NewAuthClient создает клиента и соединяется с сервером.
func NewAuthClient(addr string, logger *zap.Logger) (*AuthClient, error) {
	// Приоритет: сначала аргумент addr, затем переменная окружения AUTH_SERVICE_ADDR, затем дефолт "localhost:50051"
	if addr == "" {
		addr = os.Getenv("AUTH_SERVICE_ADDR")
		if addr == "" {
			addr = "localhost:50051"
		}
	}
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}
	c := auth.NewAuthServiceClient(conn)
	return &AuthClient{client: c, conn: conn, logger: logger}, nil
}

// Close закрывает соединение.
func (a *AuthClient) Close() error {
	return a.conn.Close()
}

// VerifyToken вызывает VerifyToken сервиса AuthService.
func (a *AuthClient) VerifyToken(ctx context.Context, token string) (*auth.VerifyTokenResponse, error) {
	req := &auth.VerifyTokenRequest{Token: token}
	resp, err := a.client.VerifyToken(ctx, req)
	if err != nil {
		a.logger.Error("Auth VerifyToken failed", zap.Error(err))
		return nil, err
	}
	return resp, nil
}
