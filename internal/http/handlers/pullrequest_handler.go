package handlers

import (
	"avito-intern-test-task-2025/internal/usecase"
	"github.com/gin-gonic/gin"
)

type PrHandler struct {
	service *usecase.PullRequestUsecase
}

func NewPrHandler(pr *usecase.PullRequestUsecase) *PrHandler {
	return &PrHandler{
		service: pr,
	}
}

func PullRequestsMerge() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func PullRequestsCreate() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func PullRequestsReassign() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
