package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

type HelloController struct {
	beego.Controller
}

func (this *HelloController) Get() {
	this.Data["title"] = "我的第一个beego"
	this.TplName = "hello.tpl"
}

//func (this *HelloController) Get() {
//	this.Data["json"] = map[string]interface{}{
//		"key":    "rex",
//		"age":    20,
//		"height": 170.5,
//	}
//
//	this.ServeJSON()
//}
