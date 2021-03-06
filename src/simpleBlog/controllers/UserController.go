package controllers

import (
	"simpleBlog/models/class"
	"github.com/astaxie/beego/validation"
	"time"
	"fmt"
	"path"
	"strings"
	"github.com/russross/blackfriday"
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
 * 加载个人主页数据
 */
func (c *UserController) Profile() {

	//获取来自页面的用户id
	username := c.Ctx.Input.Param(":username")
	//已经登录账户的访问加载对应id的
	user.UserName = username
	//通过数据库查询该用户信息
	user.Query()
	//将查询到的用户信息，存储到map中
	c.Data["u"] = user
	//查询和该用户相关的文章
	as := class.Article{}.QueryAllArticle()

	//将文章对应的tag建立关联
	for i := range as {

		if len(as[i].Content) > 0 {

			if len(as[i].Content) > 200 {
				as[i].Content = string(blackfriday.MarkdownCommon([]byte([]byte(as[i].Content)[:200])))
			} else {
				as[i].Content = string(blackfriday.MarkdownCommon([]byte(as[i].Content)))
			}

		}
	}

	if len(as) > 0 {
		//文章列表数据保存到map中
		c.Data["articles"] = as
	}

	//设置模板导向
	c.TplName = "user/profile.html"

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
	describe := c.GetString("describe")

	email := c.GetString("email")
	pwd1 := c.GetString("password")
	pwd2 := c.GetString("password2")

	valid := validation.Validation{}

	valid.Required(username, "Username")

	//验证是否传入以下参数
	valid.Required(pwd1, "Password")
	valid.Required(pwd2, "Password2")

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
		user.Describe = describe
		user.Email = email
		user.Password = class.PwGen(pwd1)

		user.RegTime = time.Now()

		switch {
		case user.ExistId():
			valid.Error("用户名被占用")
		case user.ExistEmail():
			valid.Error("邮箱被占用")
		default:
			err := user.Insert()
			if err == nil {
				c.DoLogin(user)
				c.ret.Ok = true
				c.ret.Content = user.UserName
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
 * 用户登录
 */
func (c *UserController) Login() {

	defer func() {
		c.Data["json"] = c.ret
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
	case u.Query() != nil:
		valid.Error("用户不存在")
	case class.PwCheck(pwd, u.Password) == false:
		valid.Error("密码错误")
	default:
		c.DoLogin(*u)
		c.ret.Ok = true
		c.ret.Content = u.UserName
		return
	}

	c.ret.Content = valid.Errors[0].Key + valid.Errors[0].Message
	c.ret.Ok = false
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
	user.Email = c.GetString("email")
	user.Describe = c.GetString("describe")
	user.About = c.GetString("hobby")

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
 * 替换或上传头像
 */
func (c *UserController) ResetUserLogoImg() {

	f, h, _ := c.GetFile("imgFiles") //获取上传的文件

	filename := h.Filename

	f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况

	c.SaveToFile("imgFiles", path.Join("static/img", filename))

	var err error

	/*替换图片*/
	if len(user.LogoImgUrl) > 0 {
		rest := strings.Split(user.LogoImgUrl, qiuNiuUrl)
		err = class.CoversimeUploadFile(rest[1], filename)
	} else {
		//上传图片
		err = class.SimeUploadFile(filename)
	}

	defer func() {
		c.Data["json"] = c.ret
		c.ServeJSON()
	}()

	if err == nil {
		c.Data["img"] = qiuNiuUrl + filename
		user := c.GetSession("user").(class.User)
		user.LogoImgUrl = qiuNiuUrl + filename
		user.Update()
		c.DoLogin(user)
		c.Redirect("/", 302)
	} else {
		c.ret.Ok = false
		c.ret.Content = err
		c.Data["json"] = c.ret
		c.ServeJSON()
	}

}
