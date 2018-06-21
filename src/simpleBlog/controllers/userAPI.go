package controllers

import (
	"github.com/astaxie/beego/validation"
	"simpleBlog/models/class"
	"time"
	"fmt"
	"strconv"
	"crypto/sha1"
	"crypto/md5"
	"encoding/base64"
)

func (c *UserController) API_Profile() {

	type user struct {
		Userid string
		Hobby  []string
	}

	u := user{
		"jike",
		[]string{"chess", "code"},
	}

	c.Data["json"] = u
	c.ServeJSON()

}

type RET struct {
	Ok      bool        `json:"success"`
	Content interface{} `json:"content"`
}

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

	if len(nick) < 1 {
		nick = id
	}

	valid := validation.Validation{}

	valid.Email(email, "Email")
	valid.Required(id, "Userid")
	valid.Required(pwd1, "Password")
	valid.Required(pwd2, "Password2")

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

func (c *UserController) Logout() {
	c.DoLogout()
}

func (c *UserController) Setting() {
	c.CheckLogin()

	switch c.GetString("do") {
	case "info":
		c.SettingInfo()
	case "chpwd":
		c.SettingPwd()

	}
}

func (c *UserController) SettingInfo() {
	user := c.GetSession("user").(class.User)

	user.Nick = c.GetString("nick")
	user.Email = c.GetString("email")
	user.Url = c.GetString("url")
	user.Hobby = c.GetString("hobby")

	user.Update()
	c.DoLogin(user)

	ret := RET{
		Ok: true,
	}

	c.Data["json"] = ret
	c.ServeJSON()
}

func (c *UserController) SettingPwd() {
	user := c.GetSession("user").(class.User)
	user.Password = PwGen(c.GetString("pwd2"))
	user.Update()
	c.DoLogin(user)

	ret := RET{
		Ok: true,
	}

	c.Data["json"] = ret

	c.ServeJSON()

}

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

	println(u.Password)
	println(pwd)

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

func PwCheck(pwd, saved string) bool {

	saved = Base64Decode(saved)

	if len(saved) < 4 {
		return false
	}

	salt := saved[len(saved)-4:]

	return Sha1(Md5(pwd)+salt)+salt == saved
}

func PwGen(pass string) string {
	salt := strconv.FormatInt(time.Now().UnixNano()%9000+1000, 10)
	return Base64Encode(Sha1(Md5(pass)+salt) + salt)
}

func Sha1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64Decode(s string) string {
	res, _ := base64.StdEncoding.DecodeString(s)
	return string(res)
}
