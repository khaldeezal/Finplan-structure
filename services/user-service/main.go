package main

import (
	"database/sql"
	"net"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/user"
	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/delivery"
	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/repo"
	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/service"
)

func main() {
	// Инициализируем zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// Подключаемся к PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	addr := os.Getenv("SERVICE_ADDR")
	if addr == "" {
		addr = ":50053"
	}

	// Открываем TCP-порт для grpc-сервера
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to open port", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	// Инициализируем зависимости
	userRepo := repo.NewPostgresUserRepository(db, logger)
	userService := service.NewUserService(userRepo, logger)
	userHandler := delivery.NewUserHandler(userService)

	userpb.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer) // включаем reflection для grpcurl и отладки

	logger.Info("✅ UserService started", zap.String("addr", addr))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("gRPC server terminated with error", zap.Error(err))
	}
}
