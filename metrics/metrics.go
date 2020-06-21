package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
	"log"
)

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

// AdmissionLatency measures latency / execution time of Admission Control execution
// usual usage pattern is: timer := NewAdmissionLatency() ; compute ; timer.Observe()
type RequestLatency struct {
	histo *prometheus.HistogramVec
	start time.Time
}

func Register() {
	//多增加一个接口
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(requestResource)
}

// NewAdmissionLatency provides a timer for admission latency; call Observe() on it to measure
func NewAdmissionLatency() *RequestLatency {
	return &RequestLatency{
		histo: requestLatency,
		start: time.Now(),
	}
}

// Observe measures the execution time from when the AdmissionLatency was created
func (t *RequestLatency) Observe() {
	(*t.histo).WithLabelValues().Observe(time.Now().Sub(t.start).Seconds())
}

// RequestIncrease increases the counter of request handled by this service
func RequestIncrease() {
	requestCount.WithLabelValues().Add(1)
}

func RequestResourceUpdate(rate float64)  {
	requestResource.WithLabelValues().Set(rate)
	// 输出利用率的信息
	log.Println("rate="+strconv.FormatFloat(rate,'E',-1,64))
}
