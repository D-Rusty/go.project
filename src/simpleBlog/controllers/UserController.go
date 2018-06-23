package controllers

import (
	"simpleBlog/models/class"
	"github.com/astaxie/beego/validation"
	"time"
	"fmt"
)

type UserController struct {
	BaseController
	ret RET
}

/**
 * 个人主页
 */
func (c *UserController) Profile() {
	class.InitData()
	//获取来自页面的用户id
	username := c.Ctx.Input.Param(":username")
	u := class.User{}

	if len(username) <= 0 {
		u.UserName = "drusty"
	} else {
		//已经登录账户的访问加载对应id的
		u.UserName = username
	}

	//通过数据库查询该用户信息
	u.ReadDB()
	//将查询到的用户信息，存储到map中
	c.Data["u"] = u
	//查询和该用户相关的文章
	as := class.Article{Author: &u}.Gets()
	//查询和该用户相关文章的评论
	replys := class.Reply{Author: &u}.Gets()
	//文章列表数据保存到map中
	c.Data["articles"] = as
	//文章评论数据存储到map中
	c.Data["replys"] = replys
	//设置模板导向
	c.TplName = "user/profile.html"
}

/**
 * 未登录前主页
 */
func (c *UserController) UnLoginHomePage() {
	c.TplName = "user/homepage.html"
}

/**
 * 注册页面
 */
func (c *UserController) RegisterPage() {
	c.TplName = "user/register.html"
}

/**
 * 用户注册
 */
func (c *UserController) Register() {

	defer func() {
		c.Data["json"] = c.ret
		c.ServeJSON()
	}()

	username := c.GetString("username")
	nick := c.GetString("nick")
	logoImgUrl := c.GetString("logoImgUrl")
	describe := c.GetString("describe")
	hobby := c.GetString("hobby")
	email := c.GetString("email")
	pwd1 := c.GetString("password")
	pwd2 := c.GetString("password2")

	valid := validation.Validation{}

	valid.Required(username, "Username")
	valid.Required(nick, "Nick")

	//验证是否传入以下参数
	valid.Required(pwd1, "Password")
	valid.Required(pwd2, "Password2")

	//验证别名长度
	valid.MaxSize(nick, 30, "Nick")
	valid.MinSize(nick, 3, "Nick")

	valid.MaxSize(username, 30, "Username")
	valid.MinSize(username, 3, "Username")

	valid.MaxSize(pwd1, 30, "password")
	valid.MinSize(pwd1, 3, "password")

	valid.MaxSize(pwd2, 10, "password2")
	valid.MinSize(pwd2, 6, "password2")

	//验证email格式是否正确
	valid.Email(email, "Email")

	switch {

	case valid.HasErrors():

	case pwd1 != pwd2:
		valid.Error("两次密码不一致")

	default:
		u := &class.User{
			UId:        class.GenerateRandomUserId(),
			UserName:   username,
			Nick:       nick,
			LogoImgUrl: logoImgUrl,
			Describe:   describe,
			Hobby:      hobby,
			Email:      email,
			Password:   class.PwGen(pwd1),
			PostNum:    0,
			TagNum:     0,
			RegTime:    time.Now(),
		}
		switch {
		case u.ExistId():
			valid.Error("用户名被占用")
		case u.ExistEmail():
			valid.Error("邮箱被占用")
		default:
			err := u.CreateDB()
			if err == nil {
				c.ret.Ok = true
				c.ret.Content = "注册成功"
				return
			} else {
				valid.Error(fmt.Sprintf("%v", err))
			}

		}

	}

	c.ret.Ok = false

	c.ret.Content = valid.Errors[0].Key + valid.Errors[0].Message

	return
}

/**
 * 用户设置
 */
func (c *UserController) UserSetting() {
	//检查用户是否已经登录
	c.CheckLogin()
	c.TplName = "user/setting.html"
}

/*
 * 提交设置信息变更请求
 */
func (c *UserController) Setting() {
	c.CheckLogin()

	switch c.GetString("do") {
	case "info":
		c.SettingInfo()
	case "chpwd":
		c.SettingPwd()

	}
}

/**
 * 登录页面
 */
func (c *UserController) LoginPage() {
	c.TplName = "user/login.html"
}

/**
 * 用户登录
 */
func (c *UserController) Login() {

	ret := RET{
		Ok:      true,
		Content: "success",
	}

	defer func() {
		c.Data["json"] = ret
		c.ServeJSON()
	}()

	username := c.GetString("username")
	pwd := c.GetString("password")

	valid := validation.Validation{}

	valid.Required(username, "Username")
	valid.Required(pwd, "Password")

	u := &class.User{UserName: username}

	switch {
	case valid.HasErrors():
	case u.ReadDB() != nil:
		valid.Error("用户不存在")
	case class.PwCheck(pwd, u.Password) == false:
		valid.Error("密码错误")
	default:
		c.DoLogin(*u)
		ret.Ok = true
		ret.Content = u.UserName
		return
	}

	ret.Content = valid.Errors[0].Key + valid.Errors[0].Message
	ret.Ok = false
	return
}

/**
 *退出登录
 */
func (c *UserController) Logout() {
	c.DoLogout()
}

/**
 * 基本信息设置
 */
func (c *UserController) SettingInfo() {

	user := c.GetSession("user").(class.User)
	user.Nick = c.GetString("nick")
	user.Email = c.GetString("email")
	user.Hobby = c.GetString("hobby")
	err := user.Update()

	if err == nil {
		c.DoLogin(user)
		c.ret.Ok = true
		c.ret.Content = "用户信息设置成功"
	} else {
		c.ret.Ok = false
		c.ret.Content = "用户信息设置失败"
	}

	defer func() {

		c.Data["json"] = c.ret

		c.ServeJSON()
	}()

}

/**
 * 修改密码
 */
func (c *UserController) SettingPwd() {

	user := c.GetSession("user").(class.User)

	if class.PwCheck(c.GetString("pwd1"), user.Password) == false {
		c.ret.Ok = false
		c.ret.Content = "密码输入错误"
	} else {
		err := user.Update()

		c.DoLogin(user)

		if err == nil {
			c.ret.Ok = true
			c.ret.Content = "密码设置成功"
		} else {
			c.ret.Ok = false
			c.ret.Content = "密码设置失败"
		}
	}

	defer func() {
		c.Data["json"] = c.ret
		c.ServeJSON()
	}()

}
