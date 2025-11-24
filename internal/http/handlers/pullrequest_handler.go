package handlers

import (
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PrHandler struct {
	service *usecase.PullRequestUsecase
}

func NewPrHandler(pr *usecase.PullRequestUsecase) *PrHandler {
	return &PrHandler{
		service: pr,
	}
}

func (ph PrHandler) PullRequestsMerge() gin.HandlerFunc {
	return func(c *gin.Context) {

		var pr dto.PullRequestDTO
		if err := c.BindJSON(&pr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if !c.Request.URL.Query().Has("team_name") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required param team_name",
			})
			return
		}

		r := ph.service.Merge(pr)
		c.JSON(http.StatusOK, r)
	}
}

func (ph PrHandler) PullRequestsCreate() gin.HandlerFunc {
	return func(c *gin.Context) {

		var pr dto.PullRequestDTO
		if err := c.BindJSON(&pr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if !c.Request.URL.Query().Has("team_name") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required param team_name",
			})
			return
		}

		r := ph.service.Create(pr)
		c.JSON(http.StatusOK, r)
	}
}

func (ph PrHandler) PullRequestsReassign() gin.HandlerFunc {
	return func(c *gin.Context) {

		var pr dto.PullRequestReassignDTO
		if err := c.BindJSON(&pr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if !c.Request.URL.Query().Has("team_name") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required param team_name",
			})
			return
		}

		r := ph.service.Reassign(pr)
		c.JSON(http.StatusOK, r)
	}
}
