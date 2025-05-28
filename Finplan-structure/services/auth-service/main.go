package main

import (
	"net"
	"os"

	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/handlers"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/repo"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/services"
	"go.uber.org/zap"

	"github.com/joho/godotenv"
	authpb "github.com/khaldeezal/Finplan-structure/proto-definitions/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error()) // TODO: Подумать над graceful shutdown вместо panic
	}
	defer logger.Sync()

	// Загружаем переменные окружения из .env
	if err := godotenv.Load(); err != nil {
		logger.Error("Ошибка загрузки .env файла", zap.Error(err)) // TODO: Можно добавить алертинг на сбои env
		os.Exit(1)
	}

	// Подключаемся к БД
	db, err := repo.ConnectDB(logger)
	if err != nil {
		logger.Error("Не удалось подключиться к БД", zap.Error(err))
		os.Exit(1)
	}

	userRepo := repo.NewPostgresUserRepository(db, logger)
	userService := services.NewAuthService(userRepo, "your-secret-key", logger)

	// Создаём TCP listener на порту 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen", zap.Error(err))
		os.Exit(1)
	}

	// Создаём новый grpc-сервер
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Регистрируем нашу реализацию AuthService на grpc-сервере
	authpb.RegisterAuthServiceServer(grpcServer, handlers.NewAuthHandler(userService, logger))

	logger.Info("AuthService grpc server is running on port :50051") // TODO: Можно выводить инфу о режиме (prod/dev)

	// Запускаем grpc-сервер и начинаем принимать запросы
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve", zap.Error(err))
		os.Exit(1)
	}

}
