package routers

import (
	"simpleBlog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//用户管理
	//首页
	beego.Router("/", &controllers.UserController{}, `get:Profile`)
	//用户登录
	beego.Router("/login", &controllers.UserController{}, `get:LoginPage;post:Login`)
	//用户注册
	beego.Router("/register", &controllers.UserController{}, `get:RegisterPage;post:Register`)
	//退出登录
	beego.Router("/logout", &controllers.UserController{}, `get:Logout`)
	//用户信息设置
	beego.Router("/setting", &controllers.UserController{}, `get:UserSetting;post:Setting`)
	//个人主页
	beego.Router("/user/:username", &controllers.UserController{}, `get:Profile`)

	//文章管理
	//新建文章
	beego.Router("/article/new", &controllers.ArticleController{}, `get:OnCreateArticlePage;post:PostNewArtic`)
	//文章详情
	beego.Router("/article/:id([0-9]+)", &controllers.ArticleController{}, `get:GetArticleDetails`)
	//删除文章
	beego.Router("/article/del/:id([0-9]+)", &controllers.ArticleController{}, `get:DelArticle`)
	//编辑文章
	beego.Router("/article/edit/:id([0-9]+)", &controllers.ArticleController{}, `get:EditArticle;post:SubmitEditArticle`)
	//文章存档页面
	beego.Router("/article/archive", &controllers.ArticleController{}, "get:Archive")
	//创建评论
	beego.Router("/create/reply", &controllers.ReplyController{}, `post:CreateReply`)

}
