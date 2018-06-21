package main

import (
	"fmt"
	"net/http"
	json2 "encoding/json"
	"encoding/base64"
	"time"
)

func headers(w http.ResponseWriter, r *http.Request) {
	h := r.Header
	fmt.Fprint(w, h)
	fmt.Fprintln(w, "\n")

	j := r.Header["Accept-Encoding"]
	fmt.Fprint(w, j)
	fmt.Fprintln(w, "\n")

	k := r.Header.Get("Accept-Encoding")
	fmt.Fprint(w, k)
	fmt.Fprintln(w, "\n")
}

func body(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))
}

func process(w http.ResponseWriter, r *http.Request) {
	//
	//r.ParseForm()
	//fmt.Fprintln(w, r.Form)
	//MultipartForm 测试
	//r.ParseMultipartForm(1024)
	fmt.Fprintln(w, "(1)", r.FormValue("hello"))
	fmt.Fprintln(w, "(2)", r.PostFormValue("hello"))
	fmt.Fprintln(w, "(3)", r.PostForm)
	fmt.Fprintln(w, "(4)", r.MultipartForm)

}

func writerExample(w http.ResponseWriter, r *http.Request) {
	str := `<html>
<head><title>Go Web Programming</title></head>
<body><h1>Hello World</h1></body>
</html>`
	w.Write([]byte(str))
}

func writerHeaderExample(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	fmt.Fprintln(w, "No such service,try next door")
}

func headerExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://google.com")
	w.WriteHeader(302)
}

type Post struct {
	User    string
	Threads []string
}

func jsonExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := &Post{
		User:    "Sau Sheong",
		Threads: []string{"first", "second", "third"},
	}

	json, _ := json2.Marshal(post)

	w.Write(json)
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:     "first_cookie",
		Value:    "Go Web Programming",
		HttpOnly: true,
	}

	c2 := http.Cookie{
		Name:     "second_cookie",
		Value:    "Manning Publications Co",
		HttpOnly: true,
	}

	http.SetCookie(w, &c1)
	http.SetCookie(w, &c2)
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	c1, err := r.Cookie("first_cookie")
	if err != nil {
		fmt.Fprintln(w, "Cannot get the first cookie")
	}

	cs := r.Cookies()

	fmt.Fprintln(w, c1)
	fmt.Fprintln(w, cs)

}

//闪现消息
func setMessage(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello World!")
	c := http.Cookie{
		Name:  "flash",
		Value: base64.URLEncoding.EncodeToString(msg),
	}

	http.SetCookie(w, &c)

}

func showMessage(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("flash")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "No message found")
		}
	} else {
		rc := http.Cookie{
			Name:    "flash",
			MaxAge:  -1,
			Expires: time.Unix(1, 0),
		}
		http.SetCookie(w, &rc)
		val, _ := base64.URLEncoding.DecodeString(c.Value)
		fmt.Fprintln(w, string(val))
	}

}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/headers", headers)

	http.HandleFunc("/body", body)

	http.HandleFunc("/process", process)

	http.HandleFunc("/writer", writerExample)

	http.HandleFunc("/writerHeader", writerHeaderExample)

	http.HandleFunc("/redirect", headerExample)

	http.HandleFunc("/json", jsonExample)

	http.HandleFunc("/set_cookie", setCookie)

	http.HandleFunc("/get_cookie", getCookie)

	http.HandleFunc("/set_message", setMessage)
	http.HandleFunc("/show_message", showMessage)

	server.ListenAndServe()
}
