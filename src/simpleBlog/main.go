package main

import (
	_ "simpleBlog/routers"
	"github.com/astaxie/beego"
	_ "simpleBlog/models"
	"strings"
)

func init() {
	//注册模板函数
	beego.AddFuncMap("split", SplitHobby)
}

func main() {
	//打开session
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}

/**
 * Template Function 分隔用户习惯
 */
func SplitHobby(s string, sep string) []string {
	return strings.Split(s, sep)
}