package main

import (
	_ "simpleBlog/routers"
	"github.com/astaxie/beego"
	_ "simpleBlog/models"
	"simpleBlog/models/class"
	"encoding/gob"
	"strings"
)

func init() {
	gob.Register(class.User{})
	beego.AddFuncMap("split", SplitHobby)
}

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}

func SplitHobby(s string, sep string) []string {
	return strings.Split(s, sep)
}
