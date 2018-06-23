package controllers

import (
	"github.com/astaxie/beego"
	"simpleBlog/models/class"
)

type BaseController struct {
	beego.Controller
}

type RET struct {
	Ok      bool        `json:"success"`
	Content interface{} `json:"content"`
}

/**
 * 数据准备，已登录用户，直接将session中用户数据赋值给当前控制器
 */
func (c *BaseController) Prepare() {
	if c.IsLogin() {
		c.Data["user"] = c.GetSession("user").(class.User)
	}
}

/**
 *  登录成功后将用户数据保存到session中
 */
func (c *BaseController) DoLogin(u class.User) {
	c.SetSession("user", u)
}

/**
 * 退出登录，清除session，并将页面重定向到未登录的首页
 */
func (c *BaseController) DoLogout() {
	c.DestroySession()
	c.Redirect("/login", 302)
}

/*
 * 用户是否已经登录:true已登录，false未登录
 */
func (c *BaseController) IsLogin() bool {
	return c.GetSession("user") != nil
}

/**
 * 检查用户是否已经登录，如果没有登录，将页面重定向到未登录的首页(homepage.html)
 */
func (c *BaseController) CheckLogin() {
	if !c.IsLogin() {
		c.Redirect("/", 302)
		c.Abort("302")
	}
}
