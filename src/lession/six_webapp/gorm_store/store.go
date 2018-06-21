package main

import (
	"time"
	_ "github.com/bmizerany/pq"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Post struct {
	Id       int
	Content  string
	Author   string `sql:"not null"`
	Comment  []Comment
	CreateAt time.Time
}

type Comment struct {
	Id        int
	Content   string
	Author    string `sql:"not null"`
	PostId    int    `sql:"index"`
	CreatedAt time.Time
}

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("postgres", "user=postgres dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}

	Db.AutoMigrate(&Post{}, &Comment{})
}

func main() {

	post := Post{Content: "Hello World!", Author: "Sau Sheong"}
	fmt.Println(post)
	//创建帖子
	Db.Create(&post)
	fmt.Println(post)
	//添加评论
	comment := Comment{Content: "Good post!", Author: "Joe"}
	Db.Model(&post).Association("Comments").Append(comment)

	//通过帖子获取评论
	var readPost Post
	Db.Where("author = $1", "Sau Sheong").First(&readPost)

	var comments []Comment
	Db.Model(&readPost).Related(&comments)

	fmt.Println(comments[0])
}
