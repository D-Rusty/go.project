package controllers

import (
	"simpleBlog/models/class"
	"github.com/astaxie/beego/validation"
	"time"
	"fmt"
	"strconv"
	"crypto/sha1"
	"crypto/md5"
	"encoding/base64"
)

type UserController struct {
	BaseController
	RET
}

/**
 * 个人主页
 */
func (c *UserController) Profile() {

	//获取来自页面的用户id
	id := c.Ctx.Input.Param(":id")
	u := &class.User{Id: id}
	//通过数据库查询该用户信息
	u.ReadDB()
	//将查询到的用户信息，存储到map中
	c.Data["u"] = u
	//查询和该用户相关的文章
	as := class.Article{Author: u}.Gets()
	//查询和该用户相关文章的评论
	replys := class.Reply{Author: u}.Gets()
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
 * 用户注册
 */
func (c *UserController) Register() {

	ret := RET{
		Ok:      true,
		Content: "success",
	}

	defer func() {
		c.Data["json"] = ret
		c.ServeJSON()
	}()

	id := c.GetString("userid")
	nick := c.GetString("nick")
	pwd1 := c.GetString("password")
	pwd2 := c.GetString("password2")
	email := c.GetString("email")

	//当昵称长度小于1时 昵称=用户id
	if len(nick) < 1 {
		nick = id
	}

	valid := validation.Validation{}

	//验证email格式是否正确
	valid.Email(email, "Email")

	//验证是否传入以下参数
	valid.Required(id, "Userid")
	valid.Required(pwd1, "Password")
	valid.Required(pwd2, "Password2")

	//验证UserId,nick长度
	valid.MaxSize(id, 20, "UserID")
	valid.MaxSize(nick, 30, "Nick")

	switch {

	case valid.HasErrors():

	case pwd1 != pwd2:
		valid.Error("两次密码不一致")

	default:
		u := &class.User{
			Id:       id,
			Email:    email,
			Nick:     nick,
			Password: PwGen(pwd1),
			Regtime:  time.Now(),
			Private:  class.DefaultPvt,
		}
		switch {
		case u.ExistId():
			valid.Error("用户名被占用")
		case u.ExistEmail():
			valid.Error("邮箱被占用")
		default:
			err := u.CreateDB()
			if err == nil {
				return
			}
			valid.Error(fmt.Sprintf("%v", err))
		}

	}

	ret.Ok = false

	ret.Content = valid.Errors[0].Key + valid.Errors[0].Message

	return
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

	id := c.GetString("userid")
	pwd := c.GetString("password")

	valid := validation.Validation{}

	valid.Required(id, "UserId")
	valid.Required(pwd, "Password")

	u := &class.User{Id: id}

	switch {
	case valid.HasErrors():
	case u.ReadDB() != nil:
		valid.Error("用户不存在")
	case PwCheck(pwd, u.Password) == false:
		valid.Error("密码错误")
	default:
		c.DoLogin(*u)
		ret.Ok = true
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
	user.Url = c.GetString("url")
	user.Hobby = c.GetString("hobby")
	err := user.Update()

	if err == nil {
		c.DoLogin(user)
		c.RET.Ok = true
		c.RET.Content = "用户信息设置成功"
	} else {
		c.RET.Ok = false
		c.RET.Content = "用户信息设置失败"
	}

	defer func() {

		c.Data["json"] = c.RET

		c.ServeJSON()
	}()

}

/**
 * 修改密码
 */
func (c *UserController) SettingPwd() {

	user := c.GetSession("user").(class.User)

	if PwCheck(c.GetString("pwd1"), user.Password) == false {
		c.RET.Ok = false
		c.RET.Content = "密码输入错误"
	} else {
		err := user.Update()

		c.DoLogin(user)

		if err == nil {
			c.RET.Ok = true
			c.RET.Content = "密码设置成功"
		} else {
			c.RET.Ok = false
			c.RET.Content = "密码设置失败"
		}
	}

	defer func() {
		c.Data["json"] = c.RET
		c.ServeJSON()
	}()

}

/**
 * 验证登录密码是否一致
 */
func PwCheck(pwd, saved string) bool {

	saved = Base64Decode(saved)

	if len(saved) < 4 {
		return false
	}

	salt := saved[len(saved)-4:]

	return Sha1(Md5(pwd)+salt)+salt == saved
}

/**
 * 进行数据加密
 */
func PwGen(pass string) string {
	//依据时间产生一个4位长度的随机字符串
	salt := strconv.FormatInt(time.Now().UnixNano()%9000+1000, 10)
	//1.将密码明文先经过md5加密，2.将md5加密后的密文+随机字符串，3.对字符串进行sha1加密+随机字符串
	return Base64Encode(Sha1(Md5(pass)+salt) + salt)
}

/**
 * Sha1加密
 */
func Sha1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

/**
 * md5加密
 */
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

/**
 * base64 编码
 */
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

/**
 * base64 解码
 */
func Base64Decode(s string) string {
	res, _ := base64.StdEncoding.DecodeString(s)
	return string(res)
}
