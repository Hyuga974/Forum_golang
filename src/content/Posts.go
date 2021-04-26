package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

// AllPosts : kzdnzndz
func AllPosts(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	color := RandomColor()

	allCategories := "sport;anime/manga;jeux vidéos;informatique;economie;voyage;NEWS;paranormal"
	tabCategories := strings.Split(allCategories, ";")
	var tabCat []CATEGORIES
	for _, x := range tabCategories {
		oneCategorie := CATEGORIES{
			Cat:   x,
			Color: color[x],
		}
		tabCat = append(tabCat, oneCategorie)
	}

	var all_Post []POSTINFO
	var post_info POSTINFO

	db, err := sql.Open("sqlite3", "database/database.db")

	categorie := ""
	if r.Method == "POST" {
		for _, x := range tabCategories {
			if r.FormValue(x) != "" {
				categorie = x
			}
		}

	}
	post, err := db.Query("SELECT * FROM Posts ORDER BY id DESC")

	var since string
	var id int
	var user_id int
	var title string
	var body string
	var image string
	var likes int
	var comment_nb int
	var categories string
	var userinfo INFO
	for post.Next() {
		err = post.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
		CheckErr(err)
		cat := strings.Split(categories, ";")
		var tabCategories []CATEGORIES
		for _, x := range cat {
			catephemere := CATEGORIES{
				Cat:   x,
				Color: color[x],
			}
			tabCategories = append(tabCategories, catephemere)
		}
		catCheck := false
		for _, y := range tabCategories {
			if y.Cat == categorie {
				catCheck = true
				continue
			}
		}
		if catCheck == true {
			userinfo = GetUser(user_id)

			post_info = POSTINFO{
				ID:         id,
				User_ID:    user_id,
				Title:      title,
				Body:       body,
				Image:      image,
				Categories: tabCategories,
				Likes:      likes,
				Comment_Nb: comment_nb,

				Post_User_Info: userinfo,
			}
			all_Post = append(all_Post, post_info)
		}

	}
	post.Close()
	var allUsers []INFO

	users, err := db.Query("SELECT * FROM Users ORDER BY id DESC")

	var currentlyUser INFO
	var email string
	var password string
	var username string
	var description string
	var country string
	var mod int
	for users.Next() {
		err = users.Scan(&id, &username, &email, &since, &description, &password, &image, &country, &mod)
		CheckErr(err)
		currentlyUser = GetUser(id)
		allUsers = append(allUsers, currentlyUser)
	}
	users.Close()

	db.Close()

	postInfo := POSTINFO{
		AllCategories: tabCat,
	}

	data := ALLINFO{
		Self_User_Info: user,
		Post_Info:      postInfo,

		All_User:  allUsers,
		All_Posts: all_Post,
	}

	files := []string{"template/Posts.html", "template/Common.html"}

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
