package limiter

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// 方法限流
type MethodLimiter struct {
	*Limiter
}

func NewMethodLimiter() LimiterIface {
	l := &Limiter{limiterBuckets: make(map[string]*ratelimit.Bucket)}
	return &MethodLimiter{
		Limiter: l,
	}
}

// 根据context上下文获取url（去除？后的参数，并拼接Method(get\post\put\delete)）
func (li *MethodLimiter) Key(c *gin.Context) string {
	url := c.Request.RequestURI
	index := strings.Index(url, "?")
	method := c.Request.Method
	if index == -1 {
		return method + url
	}

	return method + url[:index]
}

// 根据key，获取*ratelimit.Bucket
func (li *MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	b, ok := li.limiterBuckets[key]
	return b, ok
}

// 根据rule添加规则至map中
// TODO 现在只是放入到内存中，需要添加至redis或者数据库中
func (li *MethodLimiter) AddBuckets(rules ...LimiterBucketRule) LimiterIface {
	for _, rule := range rules {
		if b, ok := li.limiterBuckets[rule.Key]; !ok {
			bucket := ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
			li.limiterBuckets[rule.Key] = bucket
		} else {
			bucket := ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
			li.limiterBuckets[rule.Key] = bucket
			fmt.Printf("该key: %s 已经存在,规则为:%+v替换为新的规则%+v\n", rule.Key, b, bucket)
		}
	}
	return li
}
