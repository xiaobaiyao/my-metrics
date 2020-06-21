package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/mem" //用于监控CPU和内存
	"log"
	"net/http"
	"os"
	"strconv"
	"example/metrics"
)

func main(){
	http.HandleFunc("/abc", index)
	http.Handle("/metrics", promhttp.Handler())
	metrics.Register()
	err := http.ListenAndServe(":5565", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

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

func Fibonacci(n int)int{
	if n<=2{
		return 1
	}else{
		return Fibonacci(n-1)+Fibonacci(n-2)
	}
}
