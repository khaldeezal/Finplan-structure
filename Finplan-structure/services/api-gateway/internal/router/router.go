package router

import (
	"github.com/gin-gonic/gin"
	"github.com/khaldeezal/Finplan-structure/services/api-gateway/internal/middleware"
)

// NewRouter инициализирует gin-роутер и навешивает маршруты для авторизации.
func NewRouter(
// Хендлеры для авторизации
	authRegisterHandler gin.HandlerFunc,
	authLoginHandler gin.HandlerFunc,
	authVerifyHandler gin.HandlerFunc,
// Хендлеры для транзакций
	transactionAddHandler gin.HandlerFunc,
	transactionListHandler gin.HandlerFunc,
	transactionDeleteHandler gin.HandlerFunc,
	transactionGetBalanceHandler gin.HandlerFunc,
// Хендлеры для пользователей
	userGetProfileHandler gin.HandlerFunc,
	userUpdateProfileHandler gin.HandlerFunc,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	api.Use(middleware.JWTMiddleware(jwtSecret))

	// Маршруты для авторизации (на эти ручки middleware не вешается)
	auth := api.Group("/auth")
	{
		// Регистрация пользователя
		auth.POST("/register", authRegisterHandler)
		// Вход пользователя
		auth.POST("/login", authLoginHandler)
		// Подтверждение пользователя (верификация)
		auth.POST("/verify", authVerifyHandler)
	}

	// Маршруты для транзакций (на эти ручки должен быть повешен middleware)
	transactions := api.Group("/transactions")
	{
		transactions.POST("", transactionAddHandler)
		transactions.GET("", transactionListHandler)
		transactions.DELETE("/:id", transactionDeleteHandler)
	}

	// Маршруты для пользователей (на эти ручки должен быть повешен middleware)
	users := api.Group("/users")
	{
		users.GET("/:id", userGetProfileHandler)    // Профиль пользователя
		users.PUT("/:id", userUpdateProfileHandler) // Обновить профиль
	}

	return r
}
