package handlers

import (
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/queries"
	"avito-intern-test-task-2025/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
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
		if !c.Request.URL.Query().Has("team_name") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required param team_name",
			})
			return
		}
		var q queries.TeamNameQuery
		c.ShouldBindQuery(&q)

		r := th.service.GetTeamMembersByName(&q)
		c.JSON(http.StatusOK, r)
	}
}

func (th TeamHandler) TeamAdd() gin.HandlerFunc {
	return func(c *gin.Context) {

		data, err := json.NewDecoder(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid body",
			})
			return
		}
		var td dto.TeamDTO
		c.ShouldBind(&td)
		r := th.service.AddTeam(&td)
		c.JSON(http.StatusOK, r)
	}
}
