package handlers

import "github.com/gin-gonic/gin"

type PrHandler struct {
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
