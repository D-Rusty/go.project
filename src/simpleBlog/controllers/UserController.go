package controllers

import (
	"simpleBlog/models/class"
	"github.com/astaxie/beego/validation"
	"time"
	"fmt"
	"path"
	"strings"
)

type UserController struct {
	BaseController
	ret RET
}

/**
 * 用户对象
 */
var user = class.User{}

const qiuNiuUrl = "http://pax6k3826.bkt.clouddn.com/"

/**
 * 个人主页
 */

func (c *UserController) Profile() {

	//获取来自页面的用户id
	username := c.Ctx.Input.Param(":username")
	//已经登录账户的访问加载对应id的
	user.UserName = username
	//通过数据库查询该用户信息
	user.ReadDB()
	//将查询到的用户信息，存储到map中
	c.Data["u"] = user
	//查询和该用户相关的文章
	as := class.Article{}.QueryAllArticle()
	//文章列表数据保存到map中
	c.Data["articles"] = as
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
		user.UId = class.GenerateRandomUserId()
		user.UserName = username
		user.Nick = nick
		user.Describe = describe
		user.Hobby = hobby
		user.Email = email
		user.Password = class.PwGen(pwd1)

		user.PostNum = 0

		user.TagNum = 0

		user.RegTime = time.Now()

		switch {
		case user.ExistId():
			valid.Error("用户名被占用")
		case user.ExistEmail():
			valid.Error("邮箱被占用")
		default:
			err := user.CreateDB()
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

	defer func() {
		c.Data["json"] = c.ret
		c.ServeJSON()
	}()

	user := c.GetSession("user").(class.User)
	user.Nick = c.GetString("nick")
	user.Email = c.GetString("email")
	user.Describe = c.GetString("describe")
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

}

/**
 * 修改密码
 */
func (c *UserController) SettingPwd() {

	defer func() {
		c.Data["json"] = c.ret
		c.ServeJSON()
	}()

	user := c.GetSession("user").(class.User)

	oldPwd := c.GetString("oldpwd")
	newPwd := c.GetString("newpwd")
	confirmpwd := c.GetString("confirmpwd")

	valid := validation.Validation{}

	switch {
	case confirmpwd != newPwd:
		valid.Error("新密码输入不一致")
	case class.PwCheck(oldPwd, user.Password) == false:
		valid.Error("密码输入错误")
	case oldPwd == newPwd && oldPwd == confirmpwd:
		valid.Error("新旧密码不可以相同")
	default:
		user.Password = class.PwGen(confirmpwd)
		err := user.Update()
		c.DoLogin(user)
		if err == nil {
			c.ret.Ok = true
			c.ret.Content = "密码设置成功"
			return
		} else {
			c.ret.Content = "密码设置失败"
		}
	}

	c.ret.Ok = false

	c.ret.Content = valid.Errors[0].Key + valid.Errors[0].Message

}

/**
 * 注册上传用户头像
 */
func (c *UserController) RegisterUserUpLoadImg() {

	f, h, _ := c.GetFile("imgFiles") //获取上传的文件

	filename := h.Filename

	f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况

	c.SaveToFile("imgFiles", path.Join("static/img", filename))

	err := class.SimeUploadFile(filename)

	if err == nil {
		//图片上传成功
		img := filename

		c.Data["img"] = qiuNiuUrl + img

		c.TplName = "user/register.html"

		user.LogoImgUrl = img

	} else {
		fmt.Println(err)
		c.ret.Ok = false
		c.ret.Content = err
		c.Data["json"] = c.ret
		c.ServeJSON()
	}

}

/**
 * 替换头像
 */

func (c *UserController) ResetUserLogoImg() {

	f, h, _ := c.GetFile("imgFiles") //获取上传的文件

	filename := h.Filename

	f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况

	c.SaveToFile("imgFiles", path.Join("static/img", filename))

	rest := strings.Split(user.LogoImgUrl, qiuNiuUrl)

	err := class.CoversimeUploadFile(rest[1], filename)

	if err != nil {
		//图片替换成功
		fmt.Println(err)
		c.ret.Ok = false
		c.ret.Content = err
		c.Data["json"] = c.ret
		c.ServeJSON()
	} else {
		c.TplName = "user/setting.html"
	}

}
