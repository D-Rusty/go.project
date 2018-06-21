package controllers

import "simpleBlog/models/class"

type UserController struct {
	BaseController
}

func (c *UserController) Profile() {
	id := c.Ctx.Input.Param(":id")
	u := &class.User{Id: id}
	u.ReadDB()

	c.Data["u"] = u

	as := class.Article{Author: u}.Gets()
	replys := class.Reply{Author: u}.Gets()

	c.Data["articles"] = as

	c.Data["replys"] = replys

	c.TplName = "user/profile.html"
}

func (c *UserController) PageJoin() {
	c.TplName = "user/join.html"
}

func (c *UserController) PageSetting() {
	c.CheckLogin()
	c.TplName = "user/setting.html"
}
