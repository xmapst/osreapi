package router

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/xmapst/osreapi/internal/docs"
	"github.com/xmapst/osreapi/internal/router/api/v1/pool"
	"github.com/xmapst/osreapi/internal/router/api/v1/status"
	"github.com/xmapst/osreapi/internal/router/api/v1/task"
	"github.com/xmapst/osreapi/internal/router/api/v1/task/step"
	"github.com/xmapst/osreapi/internal/router/api/v1/task/workspace"
	taskv2 "github.com/xmapst/osreapi/internal/router/api/v2/task"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/router/middleware/limiter"
	"github.com/xmapst/osreapi/internal/router/middleware/zap"
	"github.com/xmapst/osreapi/internal/router/types"
	"github.com/xmapst/osreapi/internal/webui"
)

// @title           OS Remote Executor API
// @version         1.0
// @description     Operating system remote execution interface.

// @contact.name   osreapi
// @contact.url    https://github.com/xmapst/osreapi/issues

// @license.name  GPL-3.0
// @license.url   https://github.com/xmapst/osreapi/blob/main/LICENSE

func New(debug bool, maxRequests int64) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(
		cors.Default(),
		func(c *gin.Context) {
			c.Header("Server", "Gin")
			c.Header("X-Server", "Gin")
			c.Header("X-Powered-By", "XMapst")
			c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
			c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
			c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			c.Header("Pragma", "no-cache")
		},
		zap.Logger,
		zap.Recovery,
	)
	if debug {
		gin.SetMode(gin.DebugMode)

		// debug pprof
		pprof.Register(router)

		// swagger docs
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// web ui
	router.Use(webui.File("/ui", "/static", "/assets", "/favicon.ico"))

	// base
	router.GET("/version", version)
	router.GET("/healthyz", healthyz)
	router.GET("/metrics", metrics)
	router.GET("/heartbeat", heartbeat)
	router.HEAD("/heartbeat", heartbeat)
	api := router.Group("/api", limiter.New(maxRequests, http.MethodPost))
	// V1
	{
		// task
		api.GET("/v1/task", task.List)
		api.POST("/v1/task", task.Post)
		api.GET("/v1/task/:task", task.Get)
		api.PUT("/v1/task/:task", task.Manager)
		// workspace
		api.GET("/v1/task/:task/workspace", workspace.Get)
		api.DELETE("/v1/task/:task/workspace", workspace.Delete)
		api.POST("/v1/task/:task/workspace", workspace.Post)
		// step
		api.GET("/v1/task/:task/step/:step", step.Get)
		api.PUT("/v1/task/:task/step/:step", step.Manager)
		// worker pool
		api.GET("/v1/pool", pool.Detail)
		api.POST("/v1/pool", pool.Post)
		// server state
		api.GET("/v1/state", status.Detail)
	}
	// V2
	{
		// task
		api.POST("/v2/task", taskv2.Post)
	}

	// endpoints
	router.Any("/api/endpoints", func(c *gin.Context) {
		render := base.Gin{Context: c}
		var res []types.Endpoint
		var scheme = "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		for _, v := range router.Routes() {
			res = append(res, types.Endpoint{
				Method: v.Method,
				Path:   fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, v.Path),
			})
		}
		render.SetRes(res)
	})

	// no method
	router.NoMethod(func(c *gin.Context) {
		render := base.Gin{Context: c}
		render.SetError(base.CodeErrNoData, errors.New("method not allowed"))
	})

	// no route
	router.NoRoute(func(c *gin.Context) {
		render := base.Gin{Context: c}
		render.SetError(base.CodeErrNoData, errors.New("the requested path does not exist"))
	})
	return router
}
