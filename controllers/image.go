package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ImageController struct {
	baseController
}

type IMAGE_POST struct {
	Image   string    `json:"image"`                  // 属性一定要大写
	Device_id string   `json:"device_id"`             // 属性一定要大写
}

//保存
func (this *ImageController) Save()  {

	//解析post json
	var image_post IMAGE_POST
	data := this.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &image_post)
	if err != nil {
		fmt.Println("json.Unmarshal is err:", err.Error())
	}
	//fmt.Println(image_post)

	//image_data64 := this.Input().Get("image")    //解析form-data发送过来的数据
	//device_id := this.Input().Get("device_id")    //解析form-data发送过来的数据
	image_data64 := image_post.Image    // raw  json  发来的数据
	device_id := image_post.Device_id   // raw  json  发来的数据

	fmt.Println("receive image from device_id ",device_id)
	image_data, err := base64.StdEncoding.DecodeString(image_data64)

	if err != nil {
		log.Fatalln(err)
	}

	// 存储文件
	image_id :=""
	save_time := strconv.FormatInt(this.begin_time,10)
	if (device_id!=""){
		image_id = device_id + "_" + save_time + ".jpg"  // 前面是时间日志，为了分层快速所以和清理

		image_path := "/file/" + this.project + "/" + time.Now().Format("2006-01-02") + "/" + string(device_id) + "/" + save_time + ".jpg"   //格式数据"2006-01-02 15:04:05.999999999 -0700 MST
		image_dir := filepath.Dir(image_path)
		fmt.Println(image_path)
		_, err := os.Stat(image_dir)
		if (err != nil || os.IsNotExist(err)){
			err :=os.MkdirAll(image_dir,os.ModePerm)
			if (err==nil){
				fmt.Println("create dir success")
			}
		}
		err = ioutil.WriteFile(image_path,image_data,0666)
	}else {
		//读取时通过vip进行了判断,因为摄像头数据要按时间进行重定向到oss
		image_id = "vip_" + save_time + ".jpg"  // 前面是时间日志，为了分层快速所以和清理
		image_path := "/file/" + this.project + "/vip/" + save_time + ".jpg"   //格式数据"2006-01-02 15:04:05.999999999 -0700 MST
		image_dir := filepath.Dir(image_path)
		fmt.Println(image_path)
		_, err := os.Stat(image_dir)
		if (err != nil || os.IsNotExist(err)){
			err :=os.MkdirAll(image_dir,os.ModePerm)
			if (err==nil){
				fmt.Println("create dir success")
			}
		}
		err = ioutil.WriteFile(image_path,image_data,0666)
	}

	//mystruct := &JSONStruct{0, "success",image_id}  //string(map[string]string{"im::image_id": image_id})
	this.Data["json"] = map[string]interface{}{"code": 0, "message": "success","data":map[string]interface{}{"im::image_id": image_id}}      // Data是输出赋值变量
	this.ServeJSON()

}


//读取
func (this *ImageController) Read()  {
	image_id := this.Input().Get("image_id")        //  获取网址里面的变量
	fmt.Println("image_id is "+image_id)
	image_path :=""
	if (strings.Contains(image_id,"vip_")){       //如果包含vip
		image_path = "/file/" + this.project + "/vip/" + strings.Replace(image_id,"vip_", "",-1)
	} else{
		if(strings.Contains(image_id,"2019-")){         //这是之前的图片
			image_path = "/file/" + this.project + "/" + strings.Join(strings.Split(image_id,"_"),"/")
		}else{
			file_name := image_id[strings.LastIndex(image_id,"_") + 1:]
			device_id := image_id[0:strings.LastIndex(image_id,"_")]
			time_image_int,_ := strconv.Atoi(file_name[0:10])
			image_date := time.Unix(int64(time_image_int),0)
			datestr := image_date.Format("2006-01-02")
			//fmt.Println(file_name,device_id,time_image_int,image_date,datestr)
			//now_date := time.Now()
			//diff := now_date.Sub(image_date)
			//fmt.Printf("difference: %d day\n" ,int(diff.Hours()/24))
			image_path = "/file/" + this.project + "/" + datestr + "/" + device_id + "/" + file_name
			//if int(diff.Hours()/24)>7{
			//	this.Redirect("/login", 302)
			//}
		}
	}
	fmt.Println(image_path)
	_, err := os.Stat(image_path)
	if (err == nil || os.IsExist(err)){
		data,err :=ioutil.ReadFile(image_path)
		if err==nil {
			this.Ctx.Output.Body(data)
		}
	}

	this.Ctx.WriteString("no image")   //直接输出字符串
}
