package controllers

import (
	"simpleBlog/models/class"
	"fmt"
	"strconv"
	"strings"
	"html/template"
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
 * 创建Tags页面页面
 */

func (c *ArticleController) GetTags() {

	tags, tagsTotal := class.Tag{}.GetAllTag()

	c.Data["tags"] = tags
	c.Data["tagsToTal"] = tagsTotal

	c.TplName = "article/tags.html"
}

func (c *ArticleController) GetTagsArticles() {

	tag := class.Tag{}.GetTagArticle(c.GetString("tagsName"))

	fmt.Println(len(tag.Articles))

	for i := range tag.Articles {
		fmt.Println(tag.Articles[i].Title)
	}

	c.Data["tag"] = tag

	c.TplName = "article/tagsdetails.html"
}

/**
 * 获取tags列表
 */

/**
 * 提交新写好的文章到服务器
 */
func (c *ArticleController) PostNewArtic() {

	c.CheckLogin()

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
		c.Data["json"] = c.ret
		c.ServeJSON()
		return
	}

	c.ret.Ok = false
	c.ret.Content = fmt.Sprint(err)
	c.Data["json"] = c.ret
	c.ServeJSON()
}

/**
 * 获取文章详细内容
 */
func (c *ArticleController) GetArticleDetails() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	a := &class.Article{Id: id}
	a.QueryArticle()
	a.Author.ReadDB()
	a.Replys = class.Reply{Article: a}.QueryAllReply()
	c.Data["article"] = a
	c.Data["replyTree"] = a.GetReplyTree()

	body := string(blackfriday.MarkdownCommon([]byte(a.Content)))

	c.Data["bodyContent"] = body
	c.TplName = "article/article.html"
}

/*
 * 删除文章
 */
func (c *ArticleController) DelArticle() {
	c.CheckLogin()
	u := c.GetSession("user").(class.User)

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	a := &class.Article{Id: id}
	a.QueryArticle()

	if u.UId != a.Author.UId {
		c.DoLogout()
	}

	a.Defunct = true
	a.Update()

	c.Redirect("/user/"+a.Author.UId, 302)
}

/**
 *进入文章编辑页面
 */
func (c *ArticleController) EditArticle() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	a := &class.Article{Id: id}
	a.QueryArticle()
	a.Author.ReadDB()
	c.Data["article"] = a
	c.TplName = "article/create.html"
}

/**
 *文章编辑页面后结果在提交
 */

func (c *ArticleController) SubmitEditArticle() {

	c.CheckLogin()

	c.ret.Ok = false
	c.ret.Content = "编辑失败"

	defer func() {
		c.Data["json"] = c.ret
		c.ServeJSON()
	}()

	u := c.GetSession("user").(class.User)

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	a := &class.Article{Id: id}

	a.QueryArticle()

	if u.UserName != a.Author.UserName {
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
	}

}

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
