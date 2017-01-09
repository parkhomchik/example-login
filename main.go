package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const indexPage = `<doctype html>
<html>
    <head>
        <title>Example Login</title>
    </head>
    <body>
        <h1>Example Login</h1>
		%s
    </body>
</html>`

const loginBlock = `<form action="/login" method="post">
    Username:<input type="text" name="username">
    Password:<input type="password" name="password">
    <input type="submit" value="Login">
</form>`

func main() {
	http.HandleFunc("/index", index)
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	//http.HandleFunc("/logoff", notFound)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Println("index page")
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(
		res,
		indexPage,
	)
}

func login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		http.Redirect(w, r, "/", 200)
	} else {
		r.ParseForm()

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		fmt.Println(len(username))

		if (username == "") || (password == "") {

			http.Redirect(w, r, "/", 301)
			return
		}

		fmt.Println("username:", username)
		fmt.Println("password:", password)

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionid", Value: sessionID(), Expires: expiration}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", 200)
	}

}

func sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
