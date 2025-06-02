package grpcclients

import (
	"context"
	"fmt"
	"time"

	transactionpb "github.com/khaldeezal/Finplan-structure/proto-definitions/gen/transaction"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type TransactionClient struct {
	client transactionpb.TransactionServiceClient
	logger *zap.Logger
}

func NewTransactionClient(conn *grpc.ClientConn, logger *zap.Logger) *TransactionClient {
	return &TransactionClient{
		client: transactionpb.NewTransactionServiceClient(conn),
		logger: logger,
	}
}

// Пример: создание транзакции
func (c *TransactionClient) AddTransaction(ctx context.Context, req *transactionpb.AddTransactionRequest) (*transactionpb.AddTransactionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.AddTransaction(ctx, req)
	if err != nil {
		c.logger.Error("AddTransaction failed", zap.Error(err))
		return nil, fmt.Errorf("AddTransaction error: %w", err)
	}
	return resp, nil
}

// Пример: получение списка транзакций
func (c *TransactionClient) ListTransactions(ctx context.Context, req *transactionpb.ListTransactionsRequest) (*transactionpb.ListTransactionsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.ListTransactions(ctx, req)
	if err != nil {
		c.logger.Error("ListTransactions failed", zap.Error(err))
		return nil, fmt.Errorf("ListTransactions error: %w", err)
	}
	return resp, nil
}

// Пример: удаление транзакции
func (c *TransactionClient) DeleteTransaction(ctx context.Context, req *transactionpb.DeleteTransactionRequest) (*transactionpb.DeleteTransactionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.DeleteTransaction(ctx, req)
	if err != nil {
		c.logger.Error("DeleteTransaction failed", zap.Error(err))
		return nil, fmt.Errorf("DeleteTransaction error: %w", err)
	}
	return resp, nil
}

// Пример: получить баланс
func (c *TransactionClient) GetBalance(ctx context.Context, req *transactionpb.GetBalanceRequest) (*transactionpb.GetBalanceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GetBalance(ctx, req)
	if err != nil {
		c.logger.Error("GetBalance failed", zap.Error(err))
		return nil, fmt.Errorf("GetBalance error: %w", err)
	}
	return resp, nil
}
