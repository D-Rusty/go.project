package class

import (
	"github.com/astaxie/beego/orm"
	"simpleBlog/models/modules"
)

type Tag struct {
	Id          int64
	Name        string     `orm:"index"`
	Articles    []*Article `orm:"reverse(many)"`
	ArticlesNum int        `orm:"-"`
}

/**
 * 查询tag
 */
func (t Tag) Get() *Tag {
	o := orm.NewOrm()
	o.QueryTable("tag").Filter("Name", t.Name).One(&t)
	if t.Id == 0 {
		return nil
	}
	return &t
}

/**
 * 读取或创建一个tag
 */
func (t Tag) GetOrNew() *Tag {
	o := orm.NewOrm()
	_, _, _ = o.ReadOrCreate(&t, "Name")
	return &t
}

var bscolor = []string{"success", "primary", "daanger", "warning"}

func (t Tag) RandColor() string {
	return bscolor[modules.RandInt(4)]
}

/**
 * 查询所有tag以及对应的文章篇数
 */

func (t Tag) GetAllTag() (tagss []Tag, total int) {

	o := orm.NewOrm()

	qs := o.QueryTable("tag")

	qs = qs.RelatedSel()

	qs.All(&tagss)

	//将文章对应的tag建立关联
	for i := range tagss {
		//将文章对应的tag数据加载到rets[i](每一篇文章中)
		_, _ = o.LoadRelated(&tagss[i], "Articles")
		tagss[i].ArticlesNum = len(tagss[i].Articles);
		total += len(tagss[i].Articles)
	}

	return tagss, total;
}

func (t Tag) GetTagArticle(tagName string) (arts Tag) {

	o := orm.NewOrm()

	o.QueryTable("tag").Filter("Name", tagName).One(&arts)

	_, _ = o.LoadRelated(&arts, "Articles")

	return

}
