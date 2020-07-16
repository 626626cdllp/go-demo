package main
//go get github.com/astaxie/beego  会下载到gopath目录下src文件夹下
import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"image_go/controllers"
)

// 定义包
var httpRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "http request count",
	},
	[]string{"method","path"},
)

var httpRequestDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "http_request_duration",
		Help: "http request duration",
	},
	[]string{"method","path"},
)



func init() {
	beego.LoadAppConfig("ini", "image_go/conf/app.conf")        // 支持 INI、XML、JSON、YAML格式，路径是相对于gopath而言的，gopath被认为是项目路径
	beego.BConfig.WebConfig.Session.SessionOn = true     //beego.BConfig保存了所有的配置，可以在代码中读取和修改
	//BeforeStatic 静态地址之前
	//BeforeRouter 寻找路由之前
	//BeforeExec 找到路由之后，开始执行相应的 Controller 之前
	//AfterExec 执行完 Controller 逻辑之后执行的过滤器
	//FinishRouter 执行完逻辑之后执行的过滤器

	// 设置过滤器，允许跨域访问
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		//AllowOrigins:      []string{"https://192.168.0.102"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	//设置过滤器，注册prometheus,最后的false表示即使执行过滤器之前已经有输出了，也要执行此过滤器
	beego.InsertFilter("*",beego.FinishRouter, func(this *context.Context) {
		httpRequestCount.WithLabelValues(this.Request.Method,this.Request.URL.Path).Inc()
		//fmt.Println("+++++++++++",this.Request.RequestURI,this.Request.URL.Path,this.ResponseWriter.Elapsed.Nanoseconds(),"++++++++++++++++")
		httpRequestDuration.WithLabelValues(this.Request.Method,this.Request.URL.Path).Observe(float64(this.ResponseWriter.Elapsed.Nanoseconds()/1000))    // 写入微秒数据
	},false)

	// 注册 prometheus 收集器
	prometheus.MustRegister(httpRequestCount)
	prometheus.MustRegister(httpRequestDuration)

	//测试函数
	beego.Get("/",func(ctx *context.Context){
		ctx.Output.Body([]byte("hello world"))
	})
	beego.Router("/:project(.*)/v1.0/image", &controllers.ImageController{}, "get:Read;post:Save")
	beego.Handler("/metrics", promhttp.Handler())    // 注册prometheus接口

}


func main() {
	fmt.Print("start server :=================================\n")
	beego.Run()
}