package main

import (
	"net"
	"os"

	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/handlers"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/repo"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/services"
	"go.uber.org/zap"

	"github.com/joho/godotenv"
	authpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}
	defer logger.Sync()

	// Переменные окружения из .env
	if err := godotenv.Load(); err != nil {
		logger.Error("Ошибка загрузки .env файла", zap.Error(err))
		os.Exit(1)
	}

	// Подключение к БД
	db, err := repo.ConnectDB(logger)
	if err != nil {
		logger.Error("Не удалось подключиться к БД", zap.Error(err))
		os.Exit(1)
	}

	// Запуск миграции
	if err := repo.RunMigrations(db, logger); err != nil {
		logger.Fatal("Ошибка применения миграций", zap.Error(err))
	}

	userRepo := repo.NewPostgresUserRepository(db, logger)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal("JWT_SECRET is not set in environment")
	}
	userService := services.NewAuthService(userRepo, jwtSecret, logger)

	// Адрес из переменной окружения SERVICE_ADDR
	addr := os.Getenv("SERVICE_ADDR")
	if addr == "" {
		addr = ":50051"
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("failed to listen", zap.Error(err))
		os.Exit(1)
	}

	// Новый grpc-сервер
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Регистрация реализацию AuthService на grpc-сервере
	authpb.RegisterAuthServiceServer(grpcServer, handlers.NewAuthHandler(userService, logger))

	logger.Info("AuthService grpc server is running on port " + addr)

	// Запуск grpc-сервера
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve", zap.Error(err))
		os.Exit(1)
	}

}
