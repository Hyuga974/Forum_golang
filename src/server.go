package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

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
	http.HandleFunc("/editpost", content.EditPost)
	http.HandleFunc("/profil", content.Profil)
	http.HandleFunc("/login", content.Login)
	http.HandleFunc("/register", content.Register)
	fmt.Println("Start... ")
	http.ListenAndServe(":8080", nil)

}

// servHome : //* Page d'acceuil
func serveHome(w http.ResponseWriter, r *http.Request) {
	user := content.GetSession(r)
	color := content.RandomColor()

	db, err := sql.Open("sqlite3", "database/database.db")
	content.CheckErr(err)

	allPosts, _ := db.Query("SELECT * FROM Posts ORDER BY id DESC LIMIT 3")

	var post content.POSTINFO
	var mostRecent []content.POSTINFO
	var since string
	var post_id int
	var title string
	var categories string
	var body string
	var user_id int
	var image string
	var likes int
	var comments_nb int
	for allPosts.Next() {
		err = allPosts.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
		content.CheckErr(err)

		cat := strings.Split(categories, ";")
		var tabCategories []content.CATEGORIES
		for _, x := range cat {
			catephemere := content.CATEGORIES{
				Cat:   x,
				Color: color[x],
			}
			tabCategories = append(tabCategories, catephemere)
		}

		tabusers, err := db.Query("SELECT * FROM Users")
		if err != nil {
			fmt.Println(err.Error())
		}
		var userinfo content.INFO
		var userAllPost []content.POSTINFO
		var userID int
		var username string
		var email string
		var since string
		var description string
		var password string
		var country string
		for tabusers.Next() {
			err = tabusers.Scan(&userID, &username, &email, &since, &description, &password, &image, &country)
			content.CheckErr(err)
			if userID == user_id {
				userAllPost = content.GetPost(userID)
				userinfo = content.INFO{
					ID:          userID,
					Email:       email,
					PassWord:    password,
					UserName:    username,
					Since:       since,
					Description: description,
					Image:       image,
					Country:     country,
					AllPosts:    userAllPost,
				}
				break
			}
		}
		tabusers.Close()

		post = content.POSTINFO{
			ID:         post_id,
			User_ID:    user_id,
			Title:      title,
			Body:       body,
			Image:      image,
			Categories: tabCategories,
			Likes:      likes,
			Comment_Nb: comments_nb,

			Post_User_Info: userinfo,
		}
		mostRecent = append(mostRecent, post)
	}

	allPosts.Close()

	allPosts, _ = db.Query("SELECT * FROM Posts ORDER BY likes DESC LIMIT 3")
	var mostLikes []content.POSTINFO
	for allPosts.Next() {
		err = allPosts.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
		content.CheckErr(err)

		cat := strings.Split(categories, ";")
		var tabCategories []content.CATEGORIES
		for _, x := range cat {
			catephemere := content.CATEGORIES{
				Cat:   x,
				Color: color[x],
			}
			tabCategories = append(tabCategories, catephemere)
		}

		tabusers, err := db.Query("SELECT * FROM Users")
		if err != nil {
			fmt.Println(err.Error())
		}
		var userinfo content.INFO
		var userAllPost []content.POSTINFO
		var userID int
		var username string
		var email string
		var since string
		var description string
		var password string
		var country string
		for tabusers.Next() {
			err = tabusers.Scan(&userID, &username, &email, &since, &description, &password, &image, &country)
			content.CheckErr(err)
			if userID == user_id {
				userAllPost = content.GetPost(userID)
				userinfo = content.INFO{
					ID:          userID,
					Email:       email,
					PassWord:    password,
					UserName:    username,
					Since:       since,
					Description: description,
					Image:       image,
					Country:     country,
					AllPosts:    userAllPost,
				}
				break
			}
		}
		tabusers.Close()

		post = content.POSTINFO{
			ID:         post_id,
			User_ID:    user_id,
			Title:      title,
			Body:       body,
			Image:      image,
			Categories: tabCategories,
			Likes:      likes,
			Comment_Nb: comments_nb,

			Post_User_Info: userinfo,
		}
		mostLikes = append(mostLikes, post)
	}

	allPosts.Close()

	db.Close()

	data := content.ALLINFO{
		User_Info: user,
		Post_Info: content.POSTINFO{},

		Post_Most_Recent: mostRecent,
		Post_Most_Likes:  mostLikes,
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
