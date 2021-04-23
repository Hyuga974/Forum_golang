package main

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	content "forum/src/content"
)

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", content.ServeHome)

	http.HandleFunc("/posts", content.AllPosts)
	http.HandleFunc("/post", content.OnePost)
	http.HandleFunc("/newpost", content.CreationPost)
	http.HandleFunc("/editpost", content.EditPost)
	http.HandleFunc("/deletepost", content.DeletePost)
	http.HandleFunc("/profil", content.Profil)
	http.HandleFunc("/login", content.Login)
	http.HandleFunc("/register", content.Register)
	fmt.Println("Start... ")
	http.ListenAndServe(":8080", nil)

}


