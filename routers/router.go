package routers

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xmapst/osreapi/docs"
	"github.com/xmapst/osreapi/handlers"
	"github.com/xmapst/osreapi/utils"
	"time"
)

func Router(debug bool) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	lm := NewRateLimiter(time.Minute, utils.MaxRequests, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-API-KEY")
		if key != "" {
			return key, nil
		}
		return "", errors.New("API key is missing")
	})
	if utils.MaxRequests > 0 {
		router.Use(lm.Middleware())
	}
	if debug {
		pprof.Register(router)
		// swagger doc
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.GET("/", handlers.List)
	router.GET("/:id", handlers.Get)
	router.POST("/", handlers.Post)
	return router
}

// RateKeyFunc limiter
type RateKeyFunc func(ctx *gin.Context) (string, error)

type RateLimiterMiddleware struct {
	fillInterval time.Duration
	capacity     int64
	rateKeyGen   RateKeyFunc
	limiters     map[string]*ratelimit.Bucket
}

func (r *RateLimiterMiddleware) get(ctx *gin.Context) (*ratelimit.Bucket, error) {
	key, err := r.rateKeyGen(ctx)

	if err != nil {
		return nil, err
	}

	if limiter, existed := r.limiters[key]; existed {
		return limiter, nil
	}

	limiter := ratelimit.NewBucketWithQuantum(r.fillInterval, r.capacity, r.capacity)
	r.limiters[key] = limiter
	return limiter, nil
}

func (r *RateLimiterMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limiter, err := r.get(ctx)
		if err != nil || limiter.TakeAvailable(1) == 0 {
			if err == nil {
				err = errors.New("too many requests")
			}
			_ = ctx.AbortWithError(429, err)
		} else {
			ctx.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
			ctx.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
			ctx.Next()
		}
	}
}

func NewRateLimiter(interval time.Duration, capacity int64, keyGen RateKeyFunc) *RateLimiterMiddleware {
	limiters := make(map[string]*ratelimit.Bucket)
	return &RateLimiterMiddleware{
		interval,
		capacity,
		keyGen,
		limiters,
	}
}
