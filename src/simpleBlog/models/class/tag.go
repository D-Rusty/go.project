package class

import (
	"github.com/astaxie/beego/orm"
	"simpleBlog/models/modules"
)

/**
 * 文章标签
 */
type Tag struct {
	Id          int64                            // 标签id
	Name        string     `orm:"index"`         // 标签名称
	Articles    []*Article `orm:"reverse(many)"` //标签下包含的文章数组
	ArticlesNum int        `orm:"-"`             //标签下包含的文章总数

	//todo ArticlesNum 可以通过计算得出不
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

/**
 * 产生随机颜色
 */
func (t Tag) RandColor() string {
	return bscolor[modules.RandInt(4)]
}

/**
 * 查询文章标签数据表，并获取该标签下包含的文章篇数
 */
func (t Tag) GetAllTag() (tagArray []Tag, total int) {

	o := orm.NewOrm()

	qs := o.QueryTable("tag")

	qs = qs.RelatedSel()

	qs.All(&tagArray)

	//将文章对应的tag建立关联
	for i := range tagArray {
		//查询该标签下对应的文章
		_, _ = o.LoadRelated(&tagArray[i], "Articles")
		tagArray[i].ArticlesNum = len(tagArray[i].Articles);
		total += len(tagArray[i].Articles)
	}

	return tagArray, total;
}

/**
 * 获取指定标签下的文章列表
 */
func (t Tag) GetTagArticle(tagName string) (arts Tag) {

	o := orm.NewOrm()

	o.QueryTable("tag").Filter("Name", tagName).One(&arts)

	_, _ = o.LoadRelated(&arts, "Articles")

	return

}
