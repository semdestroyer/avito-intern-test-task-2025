package ServiceDependencies

import (
	"avito-intern-test-task-2025/pkg/db"
	"github.com/gin-gonic/gin"
)

type ServiceDependencies struct {
	Router *gin.Engine
	DB     *db.DB
}
