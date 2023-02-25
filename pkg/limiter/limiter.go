package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// 是limiter的超类
type LimiterIface interface {
	Key(c *gin.Context) string
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBuckets(rules ...LimiterBucketRule) LimiterIface
}

// 限流的规则
type LimiterBucketRule struct {
	Key          string
	FillInterval time.Duration //填充间隔
	Capacity     int64         //桶的容量，例如：100
	Quantum      int64         //按照固定时间放N个令牌，桶的容量不会超过Capacity
}

// 保存了根据规则创建的限流ratelimit.Bucket
type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}
