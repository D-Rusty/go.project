package class

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
)

type Article struct {
	Id        int
	Title     string `orm:"size(80)"`
	Content   string `orm:"type(text)`
	Author    *User  `orm:"rel(fk)"`
	NumReplys int
	NumViews  int

	Replys []*Reply `orm:"-"`

	Tags []*Tag `orm:"rel(m2m)"`

	Time time.Time `orm:"auto_now_add;type(datetime)"`

	Defunct bool
}

func (a *Article) ReadDB() (err error) {
	o := orm.NewOrm()
	if err = o.Read(a); err != nil {
		beego.Info(err)
	}

	_, _ = o.LoadRelated(a, "tags")

	return
}

func (a Article) Create() (n int64, err error) {
	o := orm.NewOrm()
	if n, err = o.Insert(&a); err != nil {
		beego.Info(err)
	}
	return
}

func (a Article) Update() (err error) {
	o := orm.NewOrm()
	if _, err = o.Update(&a); err != nil {
		beego.Info(err)
	}

	m2m := o.QueryM2M(&a, "Tags")
	//m2m.Clear()
	//m2m.Add(a.Tags)
	old := Article{Id: a.Id}
	_, _ = o.LoadRelated(&old, "Tags")

	//insert
VI:
	for _, vi := range a.Tags {
		for _, vj := range old.Tags {
			if vi.Id == vj.Id {
				continue VI;
			}

			m2m.Add(vi)
		}
	}
	//del
VD:
	for _, vi := range old.Tags {
		for _, vj := range a.Tags {
			if vi.Id == vj.Id {
				continue VD
			}
		}
		m2m.Remove(vi)
	}

	return
}

func (a Article) Delete() (err error) {
	a.Defunct = true
	err = a.Update()
	return
}

func (a Article) Gets() (rets []Article) {
	o := orm.NewOrm()

	qs := o.QueryTable("article")

	if a.Author != nil {
		qs = qs.Filter("Author", a.Author)
	}

	if len(a.Tags) == 1 {
		qs = qs.Filter("Tags__Tag", a.Tags[0])
	}

	qs = qs.Filter("defunct", 0)

	qs = qs.RelatedSel()

	qs.All(&rets)

	for i := range rets {
		_, _ = o.LoadRelated(&rets[i], "Tags")
	}

	return

}

func (a *Article) GetReplyTree() (rets [] *ReplyTree) {

	//replys := Reply{Article: a}.Gets()
	replys := Reply{Article: a}.Gets()
	m := make(map[int]*ReplyTree)

	for _, reply := range replys {
		tr := &ReplyTree{
			Reply:  reply,
			Childs: make([]*ReplyTree, 0),
		}

		m[tr.Id] = tr

		fmt.Println(reply.Content)

		if reply.ParentId == 0 {
			rets = append(rets, tr)
		} else {
			m[reply.ParentId].Childs = append(m[reply.ParentId].Childs, tr)
		}
	}

	return
}
