package tests

import (
	"context"

	"bytes"
	"github.com/gin-gonic/gin"
	authpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/auth"
	"github.com/khaldeezal/Finplan-structure/services/api-gateway/internal/handlers"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

const bufSize = 1024 * 1024

type mockAuthServer struct {
	authpb.UnimplementedAuthServiceServer
	RegisterFn func(context.Context, *authpb.RegisterRequest) (*authpb.AuthResponse, error)
	LoginFn    func(context.Context, *authpb.LoginRequest) (*authpb.AuthResponse, error)
	VerifyFn   func(context.Context, *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error)
}

func (m *mockAuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	return m.RegisterFn(ctx, req)
}

func (m *mockAuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	return m.LoginFn(ctx, req)
}

func (m *mockAuthServer) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	return m.VerifyFn(ctx, req)
}

func startAuthTestServer(t *testing.T, srv *mockAuthServer) (*grpc.ClientConn, func()) {
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, srv)
	go server.Serve(listener)

	dialer := func(context.Context, string) (net.Conn, error) { return listener.Dial() }
	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	cleanup := func() {
		conn.Close()
		server.Stop()
	}
	return conn, cleanup
}

func setupAuthRouter(h *handlers.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.POST("/verify", h.VerifyToken)
	return r
}

func TestAuthHandler_Register_Success(t *testing.T) {
	srv := &mockAuthServer{
		RegisterFn: func(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
			require.Equal(t, "user@example.com", req.Email)
			require.Equal(t, "password", req.Password)
			require.Equal(t, "Tester", req.Name)
			return &authpb.AuthResponse{AccessToken: "jwt-token"}, nil
		},
	}
	conn, cleanup := startAuthTestServer(t, srv)
	defer cleanup()

	h := handlers.NewAuthHandler(conn, zap.NewNop())
	router := setupAuthRouter(h)

	body := `{"email":"user@example.com","password":"password","name":"Tester"}`
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "jwt-token")
}

func TestAuthHandler_Login_Success(t *testing.T) {
	srv := &mockAuthServer{
		LoginFn: func(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
			require.Equal(t, "user@example.com", req.Email)
			require.Equal(t, "password", req.Password)
			return &authpb.AuthResponse{AccessToken: "jwt-token"}, nil
		},
	}
	conn, cleanup := startAuthTestServer(t, srv)
	defer cleanup()

	h := handlers.NewAuthHandler(conn, zap.NewNop())
	router := setupAuthRouter(h)

	body := `{"email":"user@example.com","password":"password"}`
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "jwt-token")
}

func TestAuthHandler_VerifyToken_Success(t *testing.T) {
	srv := &mockAuthServer{
		VerifyFn: func(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
			require.Equal(t, "some-token", req.Token)
			return &authpb.VerifyTokenResponse{Valid: true, UserId: "42"}, nil
		},
	}
	conn, cleanup := startAuthTestServer(t, srv)
	defer cleanup()

	h := handlers.NewAuthHandler(conn, zap.NewNop())
	router := setupAuthRouter(h)

	body := `{"token":"some-token"}`
	req, _ := http.NewRequest(http.MethodPost, "/verify", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), `"valid":true`)
	require.Contains(t, resp.Body.String(), `"user_id":"42"`)
}

func TestAuthHandler_Register_BadRequest(t *testing.T) {
	h := handlers.NewAuthHandler(&grpc.ClientConn{}, zap.NewNop())
	router := setupAuthRouter(h)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

