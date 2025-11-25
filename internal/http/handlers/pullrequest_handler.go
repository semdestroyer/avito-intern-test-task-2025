package handlers

import (
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/errors"
	"avito-intern-test-task-2025/internal/usecase"
	"net/http"

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

func (ph PrHandler) PullRequestsMerge() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pr dto.PullRequestDTO
		if err := c.BindJSON(&pr); err != nil {
			errors.RespondWithError(c, http.StatusBadRequest, errors.InvalidInputError("Invalid JSON: "+err.Error()))
			return
		}

		r, err := ph.service.Merge(pr)
		if err != nil {
			switch err.Error() {
			case "PR_NOT_FOUND":
				errors.RespondWithError(c, http.StatusNotFound, errors.NotFoundError("PR not found"))
				return
			case "PR_ALREADY_MERGED":
				errors.RespondWithError(c, http.StatusOK, errors.PrMergedError("PR already merged"))
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"pr": r})
	}
}

func (ph PrHandler) PullRequestsCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pr dto.PullRequestDTO
		if err := c.BindJSON(&pr); err != nil {
			errors.RespondWithError(c, http.StatusBadRequest, errors.InvalidInputError("Invalid JSON: "+err.Error()))
			return
		}

		r, err := ph.service.Create(pr)
		if err != nil {
			switch err.Error() {
			case "PR_EXISTS":
				errors.RespondWithError(c, http.StatusConflict, errors.PrExistsError("PR id already exists"))
				return
			case "AUTHOR_NOT_FOUND":
				errors.RespondWithError(c, http.StatusNotFound, errors.NotFoundError("Author not found"))
				return
			default:
				errors.RespondWithError(c, http.StatusInternalServerError, errors.InvalidInputError("Internal server error: "+err.Error()))
				return
			}
		}
		c.JSON(http.StatusCreated, gin.H{"pr": r})
	}
}

func (ph PrHandler) PullRequestsReassign() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pr dto.PullRequestReassignDTO
		if err := c.BindJSON(&pr); err != nil {
			errors.RespondWithError(c, http.StatusBadRequest, errors.InvalidInputError("Invalid JSON: "+err.Error()))
			return
		}

		result, err := ph.service.Reassign(pr)
		if err != nil {
			switch err.Error() {
			case "PR_NOT_FOUND":
				errors.RespondWithError(c, http.StatusNotFound, errors.NotFoundError("PR not found"))
				return
			case "PR_MERGED":
				errors.RespondWithError(c, http.StatusConflict, errors.PrMergedError("cannot reassign on merged PR"))
				return
			case "REVIEWER_NOT_ASSIGNED":
				errors.RespondWithError(c, http.StatusConflict, errors.NotAssignedError("reviewer is not assigned to this PR"))
				return
			case "NO_CANDIDATE":
				errors.RespondWithError(c, http.StatusConflict, errors.NoCandidateError("no active replacement candidate in team"))
				return
			}
		}
		c.JSON(http.StatusOK, result)
	}
}
