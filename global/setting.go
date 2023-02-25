package global

import (
	"github.com/taoruicheng/blog-service/pkg/limiter"
	"github.com/taoruicheng/blog-service/pkg/logger"
	"github.com/taoruicheng/blog-service/pkg/setting"
)

var (
	ServerSetting        *setting.ServerSettingS
	AppSetting           *setting.AppSettingS
	DatabaseSetting      *setting.DatabaseSettingS
	MethodLimiterSetting limiter.LimiterIface
	Logger               *logger.Logger
)
