package class

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type User struct {
	UId        string    `orm:"pk;"`          //id 自动增长
	UserName   string                         //用户名
	Nick       string                         //别名
	LogoImgUrl string                         //个人头像
	Describe   string                         //个人描述
	Hobby      string    `orm:"null"`         //兴趣爱好
	Email      string    `orm:"unique"`       //邮箱地址
	Password   string                         //登录密码
	PostNum    int                            //文章数量
	TagNum     int                            //标签数量
	RegTime    time.Time `orm:"auto_now_add"` //用户注册时间
}

const (
	PR_live  = iota
	PR_login
	PR_post
)

const (
	DefaultPvt = 1<<3 - 1
)

func (u *User) ReadDB() (err error) {
	o := orm.NewOrm()
	err = o.Read(u)
	return
}

func (u User) CreateDB() (err error) {
	o := orm.NewOrm()
	_, err = o.Insert(&u)
	return
}

func (u User) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&u)
	return
}

func (u User) Delete() (err error) {
	err = u.Update()
	return
}

func (u User) ExistId() bool {
	o := orm.NewOrm()
	if err := o.Read(&u); err == orm.ErrNoRows {
		return false
	}
	return true
}

func (u User) ExistEmail() bool {
	o := orm.NewOrm()
	return o.QueryTable("user").Filter("Email", u.Email).Exist()
}

func (u User) Get() *User {
	o := orm.NewOrm()
	err := o.Read(&u)
	if err == orm.ErrNoRows {
		return nil
	}
	return &u
}
