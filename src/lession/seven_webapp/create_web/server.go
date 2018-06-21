package main

import (
	"net/http"
	"strconv"
	"path"
	"encoding/json"
	"database/sql"
	_ "github.com/bmizerany/pq"
	"fmt"
)

type Post struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

var Db *sql.DB

func init() {

	var err error
	Db, err = sql.Open("postgres", "user=postgres dbname=gwp password=gwp sslmode=disable")

	if err != nil {
		panic(err)
	}
}

func retrieve(id int) (post Post, err error) {
	post = Post{}

	err = Db.QueryRow("select id,content,author from  posts where id=$1", id).Scan(&post.Id, &post.Content, &post.Author)

	return
}

func (post *Post) create() (err error) {
	statement := "insert into posts (content,author) values ($1,$2) returning id"

	stmt, err := Db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

func (post *Post) update() (err error) {
	_, err = Db.Exec("update posts set content =$2,author=$3 where  id=$1", post.Id, post.Content, post.Author)
	return
}

func (post *Post) delete() (err error) {
	_, err = Db.Exec("delete from posts where  id=$1", post.Id)
	return
}

func handleGet(w http.ResponseWriter, r *http.Request) (err error) {

	id, err := strconv.Atoi(path.Base(r.URL.Path))

	if err != nil {
		return
	}

	post, err := retrieve(id)

	if err != nil {
		return
	}

	output, err := json.MarshalIndent(&post, "", "\t\t")
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	fmt.Println("")
	return
}

func handlePOST(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var post Post
	json.Unmarshal(body, &post)
	err = post.create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

func handlePUT(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, &post)
	err = post.update()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	err = post.delete()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePOST(w, r)
	case "PUT":
		err = handlePUT(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/post/", handleRequest)

	server.ListenAndServe()
}
