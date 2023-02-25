package tracer

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// 根据应用的服务名、Jaeger服务的地址创建JaegerTracer
func NewJaegerTracer(serviceName, agentHostPort string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{ //采样速率,生产环境系统性能很重要，所以对于所有的请求都开启 Trace 显然会带来比较大的压力，另外，大量的数据也会带来很大存储压力。为此，jaeger 支持设置采样速率，根据系统实际情况设置合适的采样频率。
			Type:  "const", //const，全量采集，采样率设置0,1 分别对应打开和关闭;probabilistic ，概率采集;rateLimiting ，限速采集;remote ，一种动态采集策略，根据当前系统的访问量调节采集策略
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			// 设置是否启动LoggingReporter； 请注意，如果logSpans没开启就一定没有结果
			LogSpans: true,
			// 请一定要注意，LocalAgentHostPort填错了就一定没有结果
			LocalAgentHostPort:  agentHostPort,
			BufferFlushInterval: 1 * time.Second, //刷新缓冲区的频率
		},
	}
	tracer, close, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	// 第二步：把tracer设置为全局，以后使用opentracing.GlobalTracer()开启span即可或opentracing.StartSpan()开启spana
	opentracing.SetGlobalTracer(tracer)
	return tracer, close, nil
}
