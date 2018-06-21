package main

import (
	_ "hello/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/astaxie/beego/config"
)

type User struct {
	Id       int
	Username string
	Password string
	Age      int
}

func init() {
	orm.RegisterDataBase("default", "mysql", "root:12345678@/jikexueyuan?charset=utf8", 30)
	orm.RegisterModel(new(User))
}

func main() {

	printConfig()

	o := orm.NewOrm()
	u := User{Id: 2}
	err := o.Read(&u)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u)
	}
	beego.Run()

}

func printConfig() {

	fmt.Println(beego.AppConfig.String("appname"))

	httpPort, err := beego.AppConfig.Int("httpport")

	if err == nil {
		fmt.Println(httpPort)
	} else {
		fmt.Println(err)
	}

	testPort, _ := beego.AppConfig.Int("httpport")

	fmt.Println("test: ", testPort)

	conf, err := config.NewConfig("json", "conf/test.json")

	if err == nil {
		val := conf.String("dev::appname")
		fmt.Println("val", val)
		port := conf.String("dev::port")
		fmt.Println("port", port)
	} else {
		fmt.Println(err)
	}

}
