package main

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	_ "github.com/bmizerany/pq"
)

type Post struct {
	Id         int
	Content    string
	AuthorName string `db: author`
}

var Db *sqlx.DB

func init() {

	var err error
	Db, err = sqlx.Open("postgres", "user=postgres dbname=gwp password=gwp sslmode=disable")

	if err != nil {
		panic(err)
	}
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRowx("select id,content,author from posts where id=$1", id).StructScan(&post)
	if err != nil {
		return
	}

	return
}

func (post *Post) Create() (err error) {
	err = Db.QueryRow("insert into posts (content,author) values ($1,$2) returning id", post.Content,
		post.AuthorName).Scan(&post.Id)
	return
}

func main() {
	post := Post{Content: "Hello World!", AuthorName: "Sau Sheong"}
	post.Create()
	fmt.Println(post)
}
