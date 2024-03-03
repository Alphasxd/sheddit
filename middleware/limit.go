package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

const OnceConsume = 1

// RateLimit 创建指定填充速率和容量大小的令牌桶
func RateLimit(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(ctx *gin.Context) {
		// if bucket.Take(OnceConsume) > 0 // take 返回的是还需要多少时间取到令牌
		// bucket.TakeAvailable(OnceConsume) 返回可用的令牌数
		// 没有足够的令牌 可以选择Sleep等待或者直接返回
		if bucket.TakeAvailable(OnceConsume) < OnceConsume {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "rate limit...",
			})
			ctx.Abort()
			return
		}
		// 可以取到令牌
		ctx.Next()
	}
}
