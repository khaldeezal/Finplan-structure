package router

import (
	"github.com/gin-gonic/gin"
	"github.com/khaldeezal/Finplan-structure/services/api-gateway/internal/middleware"
)

// gin-роутер, навешивает маршруты для авторизации.
func NewRouter(
// Для авторизации
	authRegisterHandler gin.HandlerFunc,
	authLoginHandler gin.HandlerFunc,
	authVerifyHandler gin.HandlerFunc,
// Транзакций
	transactionAddHandler gin.HandlerFunc,
	transactionListHandler gin.HandlerFunc,
	transactionDeleteHandler gin.HandlerFunc,
	transactionGetBalanceHandler gin.HandlerFunc,
// Юзеров
	userGetProfileHandler gin.HandlerFunc,
	userUpdateProfileHandler gin.HandlerFunc,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	// Маршруты для авторизации (без middleware)
	auth := api.Group("/auth")
	{
		// Регистрация
		auth.POST("/register", authRegisterHandler)
		// Вход
		auth.POST("/login", authLoginHandler)
		// Подтверждение
		auth.POST("/verify", authVerifyHandler)
	}

	// Маршруты для транзакций (с middleware)
	transactions := api.Group("/transactions")
	transactions.Use(middleware.JWTMiddleware(jwtSecret))
	{
		transactions.POST("", transactionAddHandler)
		transactions.GET("", transactionListHandler)
		transactions.DELETE("/:id", transactionDeleteHandler)
		transactions.GET("/balance", transactionGetBalanceHandler)
	}

	// Маршруты для юзеров (с middleware)
	users := api.Group("/users")
	users.Use(middleware.JWTMiddleware(jwtSecret))
	{
		users.GET("/:id", userGetProfileHandler)    // Профиль пользователя
		users.PUT("/:id", userUpdateProfileHandler) // Обновить профиль
	}

	return r
}
