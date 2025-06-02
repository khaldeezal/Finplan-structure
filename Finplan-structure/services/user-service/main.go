package main

import (
	"database/sql"
	"net"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userpb "github.com/khaldeezal/Finplan-structure/proto-definitions/gen/user"
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
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=rewq1234 dbname=khaldee sslmode=disable")
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Открываем TCP-порт для grpc-сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("Failed to open port", zap.Error(err))
	}

	// JWT middleware временно отключён для дебага ручек и тестирования (включить обратно по мере необходимости)
	// grpcServer := grpc.NewServer(
	// 	grpc.UnaryInterceptor(middleware.JWTMiddleware("your-secret-key")),
	// )
	grpcServer := grpc.NewServer()

	// Инициализируем зависимости
	userRepo := repo.NewPostgresUserRepository(db, logger)
	userService := service.NewUserService(userRepo, logger)
	userHandler := delivery.NewUserHandler(userService)

	userpb.RegisterUserServiceServer(grpcServer, userHandler)
	// Важно для GoLand: reflection нужен для grpcurl и отладки ручек
	reflection.Register(grpcServer) // включаем reflection для grpcurl и отладки

	logger.Info("✅ UserService started on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("gRPC server terminated with error", zap.Error(err))
	}
}
