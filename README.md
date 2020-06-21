# my-metrics
在示例程序的基础上，对main.go文件和metrics.go文件做如下修改：  
- metrics.go:在原有的监测数据的基础上增加资源使用率的一项，并增加显示内存资源利用率的函数：
```
var (
	//访问总次数
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total",
			Help: "Number of request processed by this service.",
		}, []string{},
	)
	// 访问时延
	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency_seconds",
			Help:    "Time spent in this service.",
			Buckets: []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 30.0, 60.0, 120.0, 300.0},
		}, []string{},
	)
	// 内存资源使用率
	requestResource = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
				Name: "request_resource",
				Help: "Request the usage of MEM resource",
			},[]string{},
		)
)

```
增加相关的接口和函数：
```
func Register() {
	//多增加一个接口
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(requestResource)
}
func RequestResourceUpdate(rate float64)  {
	requestResource.WithLabelValues().Set(rate)
	// 输出利用率的信息
	log.Println("rate="+strconv.FormatFloat(rate,'E',-1,64))
}
```
- main.go:在import时增加相关的包，在index函数中增加监测资源利用率的函数：
```
import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/mem" //用于监控CPU和内存
	"log"
	"net/http"
	"os"
	"strconv"
	"example/metrics"
)
func index(w http.ResponseWriter, r *http.Request) {
	// // 请求计数器加一，标记当前时间
	timer:=metrics.NewAdmissionLatency()
	metrics.RequestIncrease()
	// 调用函数获取内存使用率
	v,_ := mem.VirtualMemory()
	metrics.RequestResourceUpdate(v.UsedPercent)
	num:=os.Getenv("Num")
	// 查看环境变量是否有效
	if num==""{
		Fibonacci(10)
		_,err:=w.Write([]byte("there is no env Num. Computation successed\n"))
		if err!=nil{
			log.Println("err:"+err.Error()+" No\n")
		}
	}else{
		numInt,_:=strconv.Atoi(num)
		Fibonacci(numInt)
		_,err:=w.Write([]byte("there is env Num. Computation successed\n"))
		if err!=nil{
			log.Println("err:"+err.Error()+" Yes\n")
		}
	}
	timer.Observe()
}
```
