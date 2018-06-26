package controllers

import (
	"simpleBlog/models/class"
	"github.com/astaxie/beego/validation"
	"time"
	"fmt"
	"path"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"context"
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

/***

 */
func (c *UserController) UpLoadImg() {
	c.TplName = "user/register.html"
	f, h, _ := c.GetFile("imgFiles") //获取上传的文件
	filename := h.Filename
	fmt.Println(filename)
	f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	c.SaveToFile("imgFiles", path.Join("static/img", filename))

	simeUploadFile(filename)

}

/**
 * 简单的上传文件
 */
func simeUploadFile(filename string) {

	path := path.Join("static/img", filename)
	accessKey := "HlE45UT8wRJBPWBb4HIup2dKn33cWcBaq6Wo-jye"
	secretKey := "IqPCJAY-0Q90VX9vF7BNSg2a_uzGlVH8TwvOi_j0"

	localFile := path

	key := filename

	bucket := "drustydatarepo"

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	mac := qbox.NewMac(accessKey, secretKey)

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}

	cfg.Zone = &storage.ZoneHuanan

	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}

	putExTtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExTtra)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ret.Key, ret.Hash)
}
