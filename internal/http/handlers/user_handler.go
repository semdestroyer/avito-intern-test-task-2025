package handlers

import (
	"avito-intern-test-task-2025/internal/http/queries"
	"avito-intern-test-task-2025/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	s usecase.UserUsecase
}

func (h UserHandler) UserSetIsActive() gin.HandlerFunc {
	return func(c *gin.Context) {
		var q queries.UserQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "incorrect query",
			})
		}
		r := h.s.UserSetIsActive(&q)
		c.JSON(http.StatusOK, r)
	}
}

func UserGetReview() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
