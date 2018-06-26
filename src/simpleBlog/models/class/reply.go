package class

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Reply struct {
	Id       int
	Content  string    `orm:"type(text)"`
	Article  *Article  `orm:"rel(fk)"`
	Author   *User     `orm:"rel(fk)"`
	Time     time.Time `orm:"auto_now_add"`
	ParentId int
	Defunct  bool
}

/**
 * 插入评论
 */
func (a *Reply) Insert() (n int64, err error) {
	o := orm.NewOrm()
	if n, err = o.Insert(a); err != nil {
		beego.Info(err)
	}
	return
}

/**
 * 查询所有评论
 */
func (a Reply) QueryAllReply() (rets [] *Reply) {
	o := orm.NewOrm()
	qs := o.QueryTable("reply")

	//添加查询条件为指定文章id的评论
	if a.Article != nil {
		qs = qs.Filter("article_id", a.Article.Id)
	}

	qs = qs.Filter("defunct", 0)

	qs.All(&rets)

	for k := range rets {

		rets[k].Article.QueryArticle()

		rets[k].Author.ReadDB()
	}

	return
}

/***
 * 评论树
 */
type ReplyTree struct {
	*Reply
	Childs []*ReplyTree
}
