package class

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

//关于模型的定义请参考:https://beego.me/docs/mvc/model/models.md

/**
 * 文章列表
 */
type Article struct {
	Id      int
	Title   string `orm:"size(80)"`
	Content string `orm:"type(text)`

	Author *User `orm:"rel(fk)"` //多篇文章对应一个作者

	//文章与标签之间设置多对多关系,一篇文章对应多个标签，一个标签对应多篇文章
	Tags []*Tag `orm:"rel(m2m)"`
	//auto_now_add 第一次保存时才设置时间
	Time time.Time `orm:"auto_now_add;type(datetime)"`

	Defunct bool //为true是文章为不可公布状态，false为正在公布状态
}

/**
 * 插入新文章
 */
func (a Article) Insert() (n int64, err error) {
	o := orm.NewOrm()
	if n, err = o.Insert(&a); err != nil {
		beego.Info(err)
	}
	return
}

/**
 * 更新文章内容以及，文章对应的tag
 */
func (a Article) Update() (err error) {

	o := orm.NewOrm()
	_, err = o.Update(&a)
	if err != nil {
		beego.Info(err)
	}

	//建立多对多映射关系查询 通过该篇文章查询，文章对应的Tags数据库
	m2m := o.QueryM2M(&a, "Tags")
	old := Article{Id: a.Id}
	_, _ = o.LoadRelated(&old, "Tags")

	//insert 查询该篇文章已经有的tag,并且与刚才表单提交的tag做比较以及新增操作
VI:
	for _, vi := range a.Tags {
		for _, vj := range old.Tags {
			if vi.Id == vj.Id {
				continue VI;
			}
		}
		m2m.Add(vi)
	}
	//del 查询该篇文章已经有的tag，清除没有在表单提交中出现的tag
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

/**
 * 修改Defunct字段为true达到删除文章效果
 */
func (a Article) Delete() (err error) {
	a.Defunct = true
	err = a.Update()
	return
}

/**
 * 从数据库中查询单篇article记录，以及article相应的tag等
 */
func (a *Article) QueryArticle() (err error) {

	o := orm.NewOrm()

	if err = o.Read(a); err != nil {
		//日志打印
		beego.Info(err)
	}

	_, _ = o.LoadRelated(a, "tags")

	return
}

/**
 * 查询指定作者文章列表
 */
func (a Article) QueryFilterOptionsAuthroArticle() (rets []Article) {

	o := orm.NewOrm()

	qs := o.QueryTable("article")

	//添加文章作者作为过滤条件
	if a.Author != nil {
		qs = qs.Filter("Author", a.Author)
	}

	//默认加载defunct 为0 的文章
	qs = qs.Filter("defunct", 0)

	//开启关系查询
	qs = qs.RelatedSel()

	//最终查询的文章结果存放到rets数组中，到现在为止，还没有将文章对应的tag加入数组
	qs.All(&rets)

	//将文章对应的tag建立关联
	for i := range rets {
		//将文章对应的tag数据加载到rets[i](每一篇文章中)
		_, _ = o.LoadRelated(&rets[i], "Tags")
	}

	return

}

/**
 * 查询指定tag文章列表
 */
func (a Article) QueryFilterOptionsTagArticle() (rets []Article) {

	o := orm.NewOrm()

	qs := o.QueryTable("article")

	//添加第一个Tag作为过滤条件
	if len(a.Tags) == 1 {
		qs = qs.Filter("Tags__Tag", a.Tags[0])
	}

	//默认加载defunct 为0 的文章
	qs = qs.Filter("defunct", 0)

	//开启关系查询
	qs = qs.RelatedSel()

	//最终查询的文章结果存放到rets数组中，到现在为止，还没有将文章对应的tag加入数组
	qs.All(&rets)

	//将文章对应的tag建立关联
	for i := range rets {
		//将文章对应的tag数据加载到rets[i](每一篇文章中)
		_, _ = o.LoadRelated(&rets[i], "Tags")
	}

	return

}

/**
 * 查询所有文章
 */
func (a Article) QueryAllArticle() (rets []Article) {

	o := orm.NewOrm()

	qs := o.QueryTable("article")

	//默认加载defunct 为0 的文章
	qs = qs.Filter("defunct", 0)

	//开启关系查询
	qs = qs.RelatedSel()

	//最终查询的文章结果存放到rets数组中，到现在为止，还没有将文章对应的tag加入数组
	qs.All(&rets)

	//将文章对应的tag建立关联
	for i := range rets {
		//将文章对应的tag数据加载到rets[i](每一篇文章中)
		_, _ = o.LoadRelated(&rets[i], "Tags")
	}

	return

}
