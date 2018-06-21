package models

import (
	_ "github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"simpleBlog/models/class"
	"github.com/astaxie/beego"
	"fmt"
)

func init() {

	orm.Debug = true

	switch beego.AppConfig.String("DB::db") {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
		orm.RegisterDataBase("default", "mysql",
			fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8&loc=%s",
				beego.AppConfig.String("DB::user"),
				beego.AppConfig.String("DB::pass"),
				beego.AppConfig.String("DB::name"),
				`Asia%2FShanghai`, ))
	case "sqlite":
		orm.RegisterDriver("sqlite", orm.DRSqlite)
		orm.RegisterDataBase("default", "sqlite3", beego.AppConfig.String("DB::file"))

	}

	orm.RegisterModel(new(class.User), new(class.Article), new(class.Tag), new(class.Reply))

	orm.RunSyncdb("default", false, true)
}
