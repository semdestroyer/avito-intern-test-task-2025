package handlers

import (
	"avito-intern-test-task-2025/internal/http/queries"
	"avito-intern-test-task-2025/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type UserHandler struct {
	service *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		service: userUsecase,
	}
}

func (h UserHandler) UserSetIsActive() gin.HandlerFunc {
	return func(c *gin.Context) {
		var q queries.UserQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "incorrect query",
			})
		}
		validate := validator.New() //TODO: вынести или выообще не использовать
		err := validate.Struct(q)
		if err != nil {
			// Обработка ошибок валидации
			log.Fatal("Validation failed:", err)
		}
		r := h.service.UserSetIsActive(&q)
		c.JSON(http.StatusOK, r)
	}
}

func UserGetReview() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
