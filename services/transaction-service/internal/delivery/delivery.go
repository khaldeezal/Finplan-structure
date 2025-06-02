package delivery

import (
	"context"
	"time"

	transactionpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/transaction"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/model"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Реализует transactionpb.TransactionServiceServer
type TransactionHandler struct {
	transactionpb.UnimplementedTransactionServiceServer
	service service.TransactionService
}

// Создаёт новый gRPC handler
func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: s,
	}
}

// Обрабатывает добавление новой транзакции
func (h *TransactionHandler) AddTransaction(ctx context.Context, req *transactionpb.AddTransactionRequest) (*transactionpb.AddTransactionResponse, error) {
	tx := &model.Transaction{
		ID:          "", // генерируется в сервисе
		UserID:      req.GetUserId(),
		CategoryID:  req.GetCategoryId(),
		Type:        model.TransactionType(req.GetType().String()),
		Amount:      req.GetAmount(),
		Description: req.GetDescription(),
		Date:        req.GetDate().AsTime(),
		CreatedAt:   time.Now(),
	}

	id, err := h.service.AddTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	return &transactionpb.AddTransactionResponse{TransactionId: id}, nil
}

// Возвращает список транзакций пользователя
func (h *TransactionHandler) ListTransactions(ctx context.Context, req *transactionpb.ListTransactionsRequest) (*transactionpb.ListTransactionsResponse, error) {
	transactions, err := h.service.ListTransactions(ctx, req.GetUserId(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, err
	}

	var resp transactionpb.ListTransactionsResponse
	resp.Transactions = []*transactionpb.Transaction{}
	for _, t := range transactions {
		enumValue, ok := transactionpb.TransactionType_value[string(t.Type)]
		if !ok {
			enumValue = 0
		}
		resp.Transactions = append(resp.Transactions, &transactionpb.Transaction{
			Id:          t.ID,
			UserId:      t.UserID,
			CategoryId:  t.CategoryID,
			Type:        transactionpb.TransactionType(enumValue),
			Amount:      t.Amount,
			Description: t.Description,
			Date:        timestamppb.New(t.Date),
			CreatedAt:   timestamppb.New(t.CreatedAt),
		})
	}

	return &resp, nil
}

// Удаляет транзакцию по ID
func (h *TransactionHandler) DeleteTransaction(ctx context.Context, req *transactionpb.DeleteTransactionRequest) (*transactionpb.DeleteTransactionResponse, error) {
	err := h.service.DeleteTransaction(ctx, req.GetTransactionId(), req.GetUserId())
	if err != nil {
		return &transactionpb.DeleteTransactionResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &transactionpb.DeleteTransactionResponse{
		Success: true,
		Message: "Transaction deleted",
	}, nil
}

// Возвращает текущий баланс пользователя
func (h *TransactionHandler) GetBalance(ctx context.Context, req *transactionpb.GetBalanceRequest) (*transactionpb.GetBalanceResponse, error) {
	income, expense, balance, err := h.service.GetBalance(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &transactionpb.GetBalanceResponse{
		IncomeTotal:  income,
		ExpenseTotal: expense,
		Balance:      balance,
	}, nil
}
