package main

import (
	"net/http"
	"html/template"
	"path/filepath"
	"strings"
	"io/ioutil"
	"github.com/russross/blackfriday"
	"fmt"
)

func main() {
	http.HandleFunc("/", handlerequest)
	http.ListenAndServe(":8000", nil)

}

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
}

func handlerequest(w http.ResponseWriter, r *http.Request) {

	posts := getPosts()

	tmpl := template.Must(template.New("index.html").Funcs(
		template.FuncMap{"markDown": markDowner}).ParseFiles("index.html"))

	err := tmpl.ExecuteTemplate(w, "index.html", posts)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func getPosts() []Post {
	a := []Post{}
	files, _ := filepath.Glob("posts/*")
	for _, f := range files {
		file := strings.Replace(f, "posts/", "", -1)
		file = strings.Replace(file, ".md", "", -1)
		fileread, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[0])
		date := string(lines[1])
		summary := string(lines[2])
		body := strings.Join(lines[3:len(lines)], "\n")
		body = string(blackfriday.MarkdownCommon([]byte(body)))
		a = append(a, Post{title, date, summary, body, file})

	}
	return a
}

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
