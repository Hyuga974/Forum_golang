package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func AdminPosts(w http.ResponseWriter, r *http.Request) {

	user := GetSession(r)

	files := []string{"template/ModerationPosts.html", "template/Common.html"}
	var data ALLINFO

	if user.Admin || user.Modo {
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
		post, _ := db.Query("SELECT * FROM Posts ORDER BY id DESC")

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
			if catCheck {
				userinfo = GetUser(user_id)

				post_info = POSTINFO{
					ID:             id,
					User_ID:        user_id,
					Title:          title,
					Body:           body,
					Image:          image,
					Categories:     tabCategories,
					Likes:          likes,
					Comment_Nb:     comment_nb,
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

		data = ALLINFO{
			Self_User_Info: user,
			Post_Info:      postInfo,

			All_User:  allUsers,
			All_Posts: all_Post,
		}

	} else {
		files = []string{"template/404.html"}
		fmt.Println("Redirect")
	}
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

func AdminUser(w http.ResponseWriter, r *http.Request) {

	user := GetSession(r)

	files := []string{"template/ModerationUsers.html", "template/Common.html"}
	var data ALLINFO
	if user.Admin || user.Modo {
		db, err := sql.Open("sqlite3", "database/database.db")

		users, err := db.Query("SELECT * FROM Users ORDER BY id ASC")

		var currentlyUser INFO
		var allUsers []INFO
		var id int
		var username string
		var email string
		var since string
		var password string
		var description string
		var image string
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

		data = ALLINFO{
			Self_User_Info: user,

			All_User: allUsers,
		}

	} else {
		files = []string{"template/404.html"}
	}
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

	userid := r.FormValue("id")
	if user.UserName != "" {
		var username string
		var email string
		var since string
		var description string
		var password string
		var image string
		var country string

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)
		post, err := db.Query("SELECT * FROM Users WHERE id=" + userid)
		if err != nil {
			fmt.Println(err.Error())
		}

		CheckErr(err)
		for post.Next() {
			err = post.Scan(&userid, &username, &email, &since, &description, &password, &image, &country)
			CheckErr(err)
		}
		post.Close()
		userIDInt, _ := strconv.Atoi(userid)
		userTarget := GetUser(userIDInt)

		if (user.Admin || user.Modo) && !userTarget.Admin {

			del, _ := db.Prepare("DELETE from Users WHERE id=?")

			res, err := del.Exec(userid)
			CheckErr(err)

			_, err = res.RowsAffected()
			CheckErr(err)

			del.Close()

			http.Redirect(w, r, "/posts", 301)

		}
		db.Close()
	} else {
		http.Redirect(w, r, "/login", 301)
	}
}

func PromoteUser(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

	fmt.Println("Promotion en cours!!!")
	userid := r.FormValue("id")
	if user.UserName != "" {
		var username string
		var email string
		var since string
		var description string
		var password string
		var image string
		var country string
		var mod int

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)
		post, err := db.Query("SELECT * FROM Users WHERE id=" + userid)
		if err != nil {
			fmt.Println(err.Error())
		}

		CheckErr(err)
		for post.Next() {
			err = post.Scan(&userid, &username, &email, &since, &description, &password, &image, &country, &mod)
			CheckErr(err)
		}
		post.Close()
		userIDInt, _ := strconv.Atoi(userid)
		userTarget := GetUser(userIDInt)

		if (user.Admin || user.Modo) && !userTarget.Admin {
			user, _ := db.Prepare("UPDATE Users SET Mod=? WHERE id=" + userid)
			modo := 1
			_, err = user.Exec(modo)
			if err != nil {
				fmt.Println(err.Error())
			}
			user.Close()
			http.Redirect(w, r, "/profil?id={{ .User_Info.ID }}", 301)

		}
		db.Close()
	} else {
		http.Redirect(w, r, "/login", 301)
	}
}

func DemoteUser(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

	fmt.Println("Relegation en cours!!!")
	userid := r.FormValue("id")
	if user.UserName != "" {
		var username string
		var email string
		var since string
		var description string
		var password string
		var image string
		var country string
		var mod int

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)
		post, err := db.Query("SELECT * FROM Users WHERE id=" + userid)
		if err != nil {
			fmt.Println(err.Error())
		}

		CheckErr(err)
		for post.Next() {
			err = post.Scan(&userid, &username, &email, &since, &description, &password, &image, &country, &mod)
			CheckErr(err)
		}
		post.Close()
		userIDInt, _ := strconv.Atoi(userid)
		userTarget := GetUser(userIDInt)

		if (user.Admin || user.Modo) && (!userTarget.Admin && !userTarget.Modo) {
			user, _ := db.Prepare("UPDATE Users SET Mod=? WHERE id=" + userid)
			modo := 0
			_, err = user.Exec(modo)
			if err != nil {
				fmt.Println(err.Error())
			}
			user.Close()
			http.Redirect(w, r, "/profil?id={{ .User_Info.ID }}", 301)

		}
		db.Close()
	} else {
		http.Redirect(w, r, "/login", 301)
	}
}
