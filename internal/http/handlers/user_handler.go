package handlers

import (
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/errors"
	"avito-intern-test-task-2025/internal/http/queries"
	"avito-intern-test-task-2025/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
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
		var dto dto.UserIsActiveDTO
		if err := c.ShouldBindJSON(&dto); err != nil {
			errors.RespondWithError(c, http.StatusBadRequest, errors.InvalidInputError("Invalid JSON: "+err.Error()))
			return
		}

		query := queries.UserIsActiveQuery{
			UserId:   dto.UserId,
			IsActive: dto.IsActive,
		}

		r, err := h.service.UserSetIsActive(&query)
		if err != nil {
			errors.RespondWithError(c, http.StatusNotFound, errors.NotFoundError("User not found"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": r})
	}
}

func (h UserHandler) UserGetReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var q queries.UserIdQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			errors.RespondWithError(c, http.StatusBadRequest, errors.InvalidInputError("Invalid query parameters: "+err.Error()))
			return
		}

		r, err := h.service.UserGetReviews(&q)
		if err != nil {
			errors.RespondWithError(c, http.StatusNotFound, errors.NotFoundError("User not found"))
			return
		}
		c.JSON(http.StatusOK, r)
	}
}
