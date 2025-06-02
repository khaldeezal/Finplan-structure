package main

import (
	"log"
	"net"
	"os"

	transactionpb "github.com/finplan/Finplan-structure/proto-definitions/gen/transaction"
	"github.com/finplan/Finplan-structure/services/transaction-service/internal/delivery"
	"github.com/finplan/Finplan-structure/services/transaction-service/internal/repo"
	"github.com/finplan/Finplan-structure/services/transaction-service/internal/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to init zap logger: %v", err)
	}
	defer logger.Sync()

	// Подключение к БД
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	// Сборка зависимостей
	transactionRepo := repo.NewTransactionRepo(db, logger)
	transactionService := service.NewTransactionService(transactionRepo, logger)
	transactionHandler := delivery.NewTransactionHandler(transactionService)

	// gRPC сервер
	grpcServer := grpc.NewServer()
	transactionpb.RegisterTransactionServiceServer(grpcServer, transactionHandler)

	// Регистрация grpc reflection для поддержки grpcurl и других клиентов
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		logger.Fatal("failed to listen on port", zap.Error(err))
	}

	logger.Info("✅ transaction-service grpc server is running on port :50053")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("failed to serve gRPC", zap.Error(err))
	}
}
