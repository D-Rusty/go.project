package routers

import (
	"simpleBlog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//用户管理
	//首页
	beego.Router("/", &controllers.UserController{}, `get:UnLoginHomePage`)
	//用户登录
	beego.Router("/login", &controllers.UserController{}, `post:Login`)
	//用户注册
	beego.Router("/register", &controllers.UserController{}, `post:Register`)
	//退出登录
	beego.Router("/logout", &controllers.UserController{}, `get:Logout`)
	//用户信息设置
	beego.Router("/setting", &controllers.UserController{}, `get:UserSetting;post:Setting`)
	//个人主页
	beego.Router("/user/:id", &controllers.UserController{}, `get:Profile`)

	//文章管理

	//删除文章
	beego.Router("/article/del/:id([0-9]+)", &controllers.ArticleController{}, `get:Del`)
	//新建文章
	beego.Router("/article/new", &controllers.ArticleController{}, `get:PageNew;post:New`)
	//文章详情
	beego.Router("/article/:id([0-9]+)", &controllers.ArticleController{}, `get:Get`)
	//编辑文章
	beego.Router("/article/edit/:id([0-9]+)", &controllers.ArticleController{}, `get:PageEdit;post:Edit`)

	beego.Router("/archive", &controllers.ArticleController{}, "get:Archive")

	//创建评论
	beego.Router("/reply/new", &controllers.ReplyController{}, `post:New`)

}