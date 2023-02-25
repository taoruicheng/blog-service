package setting

import (
	"strings"
	"time"

	"github.com/taoruicheng/blog-service/pkg/limiter"
)

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize       int
	MaxPageSize           int
	LogSavePath           string
	LogFileName           string
	LogFileExt            string
	DefaultContextTimeout int
}

type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type ApiLimiterSettingS struct {
	Method       string
	Url          string
	FillInterval int64 //填充间隔
	Capacity     int64 //桶的容量，例如：100
	Quantum      int64 //按照固定时间放N个令牌，桶的容量不会超过Capacity
}

func (a ApiLimiterSettingS) ConvertToLimiterBucketRule() limiter.LimiterBucketRule {
	return limiter.LimiterBucketRule{
		Key:          strings.ToUpper(a.Method) + a.Url,
		FillInterval: time.Duration(a.FillInterval) * time.Second,
		Capacity:     a.Capacity,
		Quantum:      a.Quantum,
	}
}

// 根据k，读取字段内容到对应的v
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
