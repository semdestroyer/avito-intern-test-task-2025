package errors

import (
	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	TEAM_EXISTS   ErrorCode = "TEAM_EXISTS"
	PR_EXISTS     ErrorCode = "PR_EXISTS"
	PR_MERGED     ErrorCode = "PR_MERGED"
	NOT_ASSIGNED  ErrorCode = "NOT_ASSIGNED"
	NO_CANDIDATE  ErrorCode = "NO_CANDIDATE"
	NOT_FOUND     ErrorCode = "NOT_FOUND"
	INVALID_INPUT ErrorCode = "INVALID_INPUT"
)

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewErrorResponse(code ErrorCode, message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetails{
			Code:    code,
			Message: message,
		},
	}
}

func RespondWithError(c *gin.Context, code int, errorResponse *ErrorResponse) {
	c.JSON(code, errorResponse)
}

func TeamExistsError(message string) *ErrorResponse {
	return NewErrorResponse(TEAM_EXISTS, message)
}

func PrExistsError(message string) *ErrorResponse {
	return NewErrorResponse(PR_EXISTS, message)
}

func PrMergedError(message string) *ErrorResponse {
	return NewErrorResponse(PR_MERGED, message)
}

func NotAssignedError(message string) *ErrorResponse {
	return NewErrorResponse(NOT_ASSIGNED, message)
}

func NoCandidateError(message string) *ErrorResponse {
	return NewErrorResponse(NO_CANDIDATE, message)
}

func NotFoundError(message string) *ErrorResponse {
	return NewErrorResponse(NOT_FOUND, message)
}

func InvalidInputError(message string) *ErrorResponse {
	return NewErrorResponse(INVALID_INPUT, message)
}
