package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/taoruicheng/blog-service/docs"
	"github.com/taoruicheng/blog-service/global"
	"github.com/taoruicheng/blog-service/internal/model"
	"github.com/taoruicheng/blog-service/internal/routers"
	"github.com/taoruicheng/blog-service/pkg/limiter"
	"github.com/taoruicheng/blog-service/pkg/logger"
	"github.com/taoruicheng/blog-service/pkg/setting"
	"github.com/taoruicheng/blog-service/pkg/tracer"
	_ "go.uber.org/automaxprocs"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	buildTime    string
	buildVersion string
	gitCommitID  string
)

func init() {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>程序build信息>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Printf("build_time:%s\n", buildTime)
	fmt.Printf("buildVersion:%s\n", buildVersion)
	fmt.Printf("gitCommitID:%s\n", gitCommitID)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>程序build信息>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
	err = setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

// @title           博客系统
// @version         1.0
// @description     Go 语言编程之旅：一起用 Go 做项目
// @termsOfService  https://github.com/taoruicheng/blog-service

// @contact.name   taoruicheng
// @contact.url    https://github.com/taoruicheng/blog-service
// @contact.email  taoruicheng@126.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8000
// @BasePath  /
func main() {

	//设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	//开启后台线程
	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()
	//如下为优雅停止程序的代码（最多等待10秒给后台处理未处理完的请求）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}
	log.Println("Server exiting")
}
func setupSetting() error {
	s, err := setting.NewSetting()
	if err != nil {
		return err
	}
	//abc := &setting.ServerSettingS{}

	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	//创建 apilimite的slice
	apiLimits := make([]setting.ApiLimiterSettingS, 1)
	err = s.ReadSection("ApiLimiter", &apiLimits)
	methodLimiter := limiter.NewMethodLimiter()
	for _, limit := range apiLimits {
		bucketRule := limit.ConvertToLimiterBucketRule()
		methodLimiter.AddBuckets(bucketRule)
	}
	global.MethodLimiterSetting = methodLimiter
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}
func setupDBEngine() error {
	var err error
	//给全局global.DBEngine设置值
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,  //并且设置日志文件所允许的最大占用空间为 600MB
		MaxAge:    10,   //日志文件最大生存周期为 10 天
		LocalTime: true, //日志文件名的时间格式为本地时间
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func setupTracer() error {
	tracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = tracer
	return nil
}
