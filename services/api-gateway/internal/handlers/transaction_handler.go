package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	transactionpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/transaction"
)

type TransactionHandler struct {
	Client transactionpb.TransactionServiceClient
}

// Удаление транзакции по id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	userID := c.Query("user_id")

	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction_id path parameter is required"})
		return
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	req := &transactionpb.DeleteTransactionRequest{
		TransactionId: transactionID,
		UserId:        userID,
	}

	_, err := h.Client.DeleteTransaction(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

func NewTransactionHandler(client transactionpb.TransactionServiceClient) *TransactionHandler {
	return &TransactionHandler{Client: client}
}

// Добавление транзакции
func (h *TransactionHandler) AddTransaction(c *gin.Context) {
	var req transactionpb.AddTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	resp, err := h.Client.AddTransaction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add transaction"})
		return
	}
	c.JSON(http.StatusCreated, resp) // Всё, resp сериализуется как есть
}

// Список транзакций
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit query parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset query parameter"})
		return
	}

	req := &transactionpb.ListTransactionsRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resp, err := h.Client.ListTransactions(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list transactions"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Получение баланса пользователя по user_id из пути
func (h *TransactionHandler) GetBalance(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	req := &transactionpb.GetBalanceRequest{
		UserId: userID,
	}

	resp, err := h.Client.GetBalance(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
