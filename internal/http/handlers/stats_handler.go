package handlers

import (
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsService interface {
	GetAssignmentStats() (dto.AssignmentStatsDTO, error)
}

type StatsHandler struct {
	service StatsService
}

func NewStatsHandler(service StatsService) *StatsHandler {
	return &StatsHandler{
		service: service,
	}
}

func (h StatsHandler) GetAssignmentsStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.service.GetAssignmentStats()
		if err != nil {
			errors.RespondWithError(c, http.StatusInternalServerError, errors.InvalidInputError("failed to load stats: "+err.Error()))
			return
		}
		c.JSON(http.StatusOK, gin.H{"stats": stats})
	}
}
