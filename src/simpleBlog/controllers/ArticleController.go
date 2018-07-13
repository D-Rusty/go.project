package controllers

import (
	"simpleBlog/models/class"
	"fmt"
	"strconv"
	"strings"
	"github.com/russross/blackfriday"
)

type ArticleController struct {
	BaseController
	ret RET
}

/**
 * 创建文章页面
 */
func (c *ArticleController) OnCreateArticlePage() {
	c.CheckLogin()
	c.TplName = "article/create.html"
}

/**
 * 获取该用户下所有文章标签列表
 */
func (c *ArticleController) GetTags() {

	tags, tagsTotal := class.Tag{}.GetAllTag()

	//标签数组
	c.Data["tags"] = tags
	//标签总数
	c.Data["tagsToTal"] = tagsTotal

	c.TplName = "article/tags.html"

}

/**
 * 获取指定标签下面文章列表
 */
func (c *ArticleController) GetTagsArticles() {

	c.Data["tag"] = class.Tag{}.GetTagArticle(c.GetString("tagsName"))

	c.TplName = "article/tagsdetails.html"
}

/**
 * 提交新写好的文章
 */
func (c *ArticleController) PostNewArtic() {

	c.CheckLogin()

	defer func() {
		c.Data["json"] = c.ret
		c.ServeJSON()
	}()

	u := c.GetSession("user").(class.User)

	a := &class.Article{
		Title:   c.GetString("title"),
		Content: c.GetString("content"),
		Author:  &u,
	}

	//将新写好的文章插入数据库
	n, err := a.Insert()

	if err == nil {
		c.ret.Ok = true
		c.ret.Content = n
	} else {
		c.ret.Ok = false
		c.ret.Content = fmt.Sprint(err)
	}

}

/**
 * 获取文章详细内容
 */
func (c *ArticleController) GetArticleDetails() {

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	a := &class.Article{Id: id}

	//查询该文章id下对应的文章记录
	a.QueryArticle()
	//查询该篇文章作者信息
	a.Author.Query()

	c.Data["article"] = a

	body := string(blackfriday.MarkdownCommon([]byte(a.Content)))

	c.Data["bodyContent"] = body

	c.TplName = "article/article.html"
}

/*
 * 删除文章
 */
func (c *ArticleController) DelArticle() {

	//检查用户是否登录
	c.CheckLogin()
	//获取当前登录用户对象
	u := c.GetSession("user").(class.User)

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	a := &class.Article{Id: id}
	//查询该篇文章对应的数据库信息
	a.QueryArticle()

	//如果当前登录用户信息和文章用户信息不一致，则强制退出登录
	if u.UId != a.Author.UId {
		c.DoLogout()
	}

	//将文章状态改为隐藏
	a.Defunct = true
	//更新文章数据库信息
	a.Update()
	//重定向到当前登录用户主页
	c.Redirect("/user/"+u.UId, 302)
}

/**
 *进入文章编辑页面
 */
func (c *ArticleController) EditArticle() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	a := &class.Article{Id: id}
	a.QueryArticle()
	a.Author.Query()
	c.Data["article"] = a
	c.TplName = "article/create.html"
}

/**
 *文章编辑页面结果提交
 */
func (c *ArticleController) SubmitEditArticle() {

	c.CheckLogin()

	defer func() {
		c.Data["json"] = c.ret
		c.ServeJSON()
	}()

	u := c.GetSession("user").(class.User)

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	a := &class.Article{Id: id}

	a.QueryArticle()

	if u.UId != a.Author.UId {
		c.DoLogout()
	}

	strs := strings.Split(c.GetString("tag"), ",")

	tags := []*class.Tag{}

	//1.将传入的tag标签用逗号进行分割并去掉多余的空格
	//2.将传入的tag字符串插入到数据库
	//3.将传入的tag存放到tag map表中
	for _, v := range strs {
		tags = append(tags, class.Tag{Name: strings.TrimSpace(v)}.GetOrNew())
	}

	a.Title = c.GetString("title")
	a.Content = c.GetString("content")

	a.Tags = tags

	err := a.Update()

	if err == nil {
		c.ret.Ok = true
		c.ret.Content = "编辑成功"
	} else {
		c.ret.Ok = false
		c.ret.Content = "编辑失败"
	}

}
