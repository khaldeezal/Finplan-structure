package main

import (
	"github.com/joho/godotenv"
	"github.com/khaldeezal/Finplan-structure/services/api-gateway/internal/handlers"
	"github.com/khaldeezal/Finplan-structure/services/api-gateway/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"

	transactionpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/transaction"
	userpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/user"
	"go.uber.org/zap"
)

func main() {
	// Переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env:", err)
	}

	// gRPC адреса микросервисов
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	transactionServiceAddr := os.Getenv("TRANSACTION_SERVICE_ADDR")

	if userServiceAddr == "" || authServiceAddr == "" || transactionServiceAddr == "" {
		log.Println("Не все сервисы подключены! Проверь USER_SERVICE_ADDR, AUTH_SERVICE_ADDR, TRANSACTION_SERVICE_ADDR")
		return
	}

	// Подключение к gRPC сервисам
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to user service:", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	authConn, err := grpc.Dial(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to auth service:", err)
	}
	defer authConn.Close()

	transactionConn, err := grpc.Dial(transactionServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to transaction service:", err)
	}
	defer transactionConn.Close()
	transactionClient := transactionpb.NewTransactionServiceClient(transactionConn)

	logger, _ := zap.NewProduction()
	authHandler := handlers.NewAuthHandler(authConn, logger)
	transactionHandler := handlers.NewTransactionHandler(transactionClient)
	userHandler := handlers.NewUserHandler(userClient, logger)
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// Роутер
	r := router.NewRouter(
		authHandler.Register,
		authHandler.Login,
		authHandler.VerifyToken,
		transactionHandler.AddTransaction,
		transactionHandler.ListTransactions,
		transactionHandler.DeleteTransaction,
		transactionHandler.GetBalance,
		userHandler.GetUserProfile,
		userHandler.UpdateUserProfile,
		jwtSecret,
	)

	// Запуск сервера
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting REST API Gateway on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Println("HTTP server error:", err)
	}
}
