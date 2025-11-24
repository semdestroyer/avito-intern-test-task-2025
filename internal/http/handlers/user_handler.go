package handlers

import (
	"avito-intern-test-task-2025/internal/http/queries"
	"avito-intern-test-task-2025/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

		if !c.Request.URL.Query().Has("is_active") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required param is_active",
			})
			return
		}
		if _, err := strconv.ParseBool(c.Request.URL.Query().Get("is_active")); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "is_active is not bool",
			})
			return
		}

		if !c.Request.URL.Query().Has("user_id") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required param user_id",
			})
			return
		}
		var q queries.UserIsActiveQuery
		c.ShouldBindQuery(&q)

		r := h.service.UserSetIsActive(&q)
		c.JSON(http.StatusOK, r)
	}
}

func (h UserHandler) UserGetReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.Request.URL.Query().Has("user_id") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required param user_id",
			})
			return
		}
		var q queries.UserIdQuery
		c.ShouldBindQuery(&q)

		r := h.service.UserGetReviews(&q)
		c.JSON(http.StatusOK, r)
	}
}
