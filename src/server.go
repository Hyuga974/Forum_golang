package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	content "forum/src/content"
)

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/posts", content.AllPosts)
	http.HandleFunc("/post", content.OnePost)
	http.HandleFunc("/newpost", content.CreationPost)
	http.HandleFunc("/profil", content.Profil)
	http.HandleFunc("/login", content.Login)
	http.HandleFunc("/register", content.Register)
	fmt.Println("Start... ")
	http.ListenAndServe(":555", nil)

}

// servHome : //* Page d'acceuil
func serveHome(w http.ResponseWriter, r *http.Request) {
	user := content.GetSession(r)

	db, err := sql.Open("sqlite3", "database/database.db")
	content.CheckErr(err)

	//var resPost []content.POSTINFO

	//allPosts, _ := db.Query("SELECT * FROM Posts ORDER BY id DESC LIMIT 3")
	// for post, _ := range allPosts{
	// 		Transfère dans une variable
	// }

	//allPosts.Close()
	db.Close()

	data := content.ALLINFO{
		User_Info:      user,
		Post_Info:      content.POSTINFO{},
		Post_User_Info: user,
	}

	files := []string{"template/Home.html", "template/Common.html"}

	tmp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error: Check template", 500)
	}

	err = tmp.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}

}
