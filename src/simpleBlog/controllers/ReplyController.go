package controllers

import (
	"simpleBlog/models/class"
	"regexp"
	"strings"
	"fmt"
)

type ReplyController struct {
	BaseController
	RET
}

func (c *ReplyController) New() {

	c.CheckLogin()
	user := c.GetSession("user").(class.User)

	defer func() {
		c.Data["json"] = c.RET
		c.ServeJSON()
	}()

	article_id, _ := c.GetInt("article_id")

	reply := &class.Reply{
		Author:  &user,
		Article: &class.Article{Id: article_id},
		Content: c.GetString("content"),
	}

	if ok, _ := regexp.MatchString(`^\@\w+ `, reply.Content); ok {
		reply.ParentId, _ = c.GetInt("parent_id")
		reply.Content = strings.SplitN(reply.Content, " ", 2)[1]
		fmt.Print("reply.Content:")
		fmt.Println(reply.Content)
	}

	if len(reply.Content) < 1 {
		c.RET.Ok = false
		c.RET.Content = "评论不能为空"
		return
	}

	_, err := reply.Create()
	if err != nil {
		c.RET.Ok = false
		c.RET.Content = err.Error()
	}

	c.RET.Ok = true
}