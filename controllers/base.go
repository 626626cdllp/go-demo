package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
)


//公共处理函数
type baseController struct {
	beego.Controller
	begin_time int64        // 开始时间戳,微秒
	end_time int64          // 结束时间,微妙
	controllerName string   // 模块
	actionName     string   // 函数
	project string          // 项目
	path string             // 路径
}


//准备函数
func (this *baseController) Prepare()  {



	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])  //control的名称，也就是文件的名称
	this.actionName = strings.ToLower(actionName)        //试图函数的名称
	this.project = this.Ctx.Input.Param(":project")   // 获取项目
	if this.project==""{
		this.Data["json"] = map[string]interface{}{"code": -1, "message": "args project not exist","data":""}      // Data是输出赋值变量
		this.ServeJSON()
		this.StopRun()
	}
	this.path = this.Ctx.Request.URL.Path
	this.begin_time=int64(time.Now().UnixNano())/1000
	fmt.Println("===================",this.path,",",time.Now().Format("2006-01-02 15:04:05.999999"))
	// 设置允许跨域访问
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", this.Ctx.Request.Header.Get("Origin"))
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET,POST")

}

//响应结束函数
func (this *baseController) Finish()  {
	this.end_time=time.Now().UnixNano()/1000
	cost_time := this.end_time-this.begin_time
	fmt.Printf("finish %s %s,time cost %d us\n",this.controllerName,this.actionName,cost_time)
}


//获取用户IP地址
func (this *baseController) getClientIp() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}



// 响应
func (this * baseController) response(error_code int,message string,result map[string]string){
	this.Data["json"] = map[string]interface{}{"error": error_code, "message": message,"result":result}      // Data是输出赋值变量
	this.ServeJSON()
}



