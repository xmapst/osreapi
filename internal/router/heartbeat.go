package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func heartbeat(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
