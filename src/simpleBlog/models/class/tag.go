package class

import (
	"github.com/astaxie/beego/orm"
	"simpleBlog/models/modules"
)

type Tag struct {
	Id       int64
	Name     string     `orm:"index"`
	Articles []*Article `orm:"reverse(many)"`
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
