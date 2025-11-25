package handlers

import (
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/errors"
	"avito-intern-test-task-2025/internal/http/queries"
	"avito-intern-test-task-2025/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service *usecase.TeamUsecase
}

func NewTeamHandler(teamHandler *usecase.TeamUsecase) *TeamHandler {
	return &TeamHandler{
		service: teamHandler,
	}
}

func (th TeamHandler) TeamGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var q queries.TeamNameQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			errors.RespondWithError(c, http.StatusBadRequest, errors.InvalidInputError("Invalid query parameters: "+err.Error()))
			return
		}

		r, err := th.service.GetTeamMembersByName(&q)
		if err != nil {
			errors.RespondWithError(c, http.StatusNotFound, errors.NotFoundError("Team not found"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"team": r})
	}
}

func (th TeamHandler) TeamAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var td dto.TeamDTO
		if err := c.BindJSON(&td); err != nil {
			errors.RespondWithError(c, http.StatusBadRequest, errors.InvalidInputError("Invalid JSON: "+err.Error()))
			return
		}

		r, err := th.service.AddTeam(td)
		if err != nil {
			errors.RespondWithError(c, http.StatusBadRequest, errors.TeamExistsError("Team already exists"))
			return
		}
		c.JSON(http.StatusCreated, gin.H{"team": r})
	}
}
