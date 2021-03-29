package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func CreationPost(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

	if user.UserName != "" {
		var Post POSTINFO
		var tabCat []CATEGORIES

		//sport, anime/manga, economie, jeux vidéo, informatique, voyages, NEW, paranormal.
		allCategories := "sport;anime/manga;jeux vidéos;informatique;economie;voyage;NEWS;paranormal"
		tabCategories := strings.Split(allCategories, ";")

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)

		if r.Method == "POST" {
			fmt.Println("ici")
			datab, err := db.Prepare("INSERT INTO Posts (title, categories, body, user_id, image, likes, comment_nb) VALUES (?,?,?,?,?,?,?)")
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Server Error", 500)
			}
			title := r.FormValue("title")
			body := r.FormValue("body")
			image := r.FormValue("myFile")
			likes := 0
			comment_nb := 0

			fmt.Printf("image= %T\n", image)

			var categoriesCheck string
			for _, categorie := range tabCategories {
				if r.FormValue(categorie) != "" {
					categoriesCheck += categorie + ";"
				}
			}

			if title != "" && body != "" && categoriesCheck != "" {
				user_id := user.ID
				_, err := datab.Exec(title, categoriesCheck, body, user_id, image, likes, comment_nb)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			tabCategoriesCheck := strings.Split(categoriesCheck, ";")
			for _, x := range tabCategoriesCheck {
				oneCategorie := CATEGORIES{
					Cat: x,
				}
				tabCat = append(tabCat, oneCategorie)
			}
			Post = POSTINFO{
				User_ID:    user.ID,
				Title:      title,
				Body:       body,
				Image:      image,
				Categories: tabCat,
			}

			uId := strconv.Itoa(user.ID)

			newPost, err := db.Query("SELECT * FROM Posts WHERE user_id=" + uId + " ORDER BY id DESC LIMIT 1")
			fmt.Println(newPost)
			if err != nil {
				fmt.Println(err.Error())
			}
			var id string
			var categories string
			var user_id int
			for newPost.Next() {
				err = newPost.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb)
				CheckErr(err)
			}
			newPost.Close()
			fmt.Println("id post:  ", id)

			http.Redirect(w, r, "/post?id="+id, 301)
		} else {
			for _, x := range tabCategories {
				oneCategorie := CATEGORIES{
					Cat: x,
				}
				tabCat = append(tabCat, oneCategorie)
			}
			Post = POSTINFO{
				Categories: tabCat,
			}
		}
		data := ALLINFO{
			User_Info: user,
			Post_Info: Post,
		}

		fmt.Println(Post)
		files := []string{"template/CreatePost.html", "template/Common.html"}
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
	} else {
		fmt.Println("LA")
		http.Redirect(w, r, "/login", 301)
	}
	fmt.Println("Test posts")
}

//OnePost : Page pour un seul post
func OnePost(w http.ResponseWriter, r *http.Request) {
	userInfo := GetSession(r)

	likeNow := ""
	likes := 0

	post_id := r.FormValue("id")
	upost_id, err := strconv.Atoi(post_id)

	db, _ := sql.Open("sqlite3", "database/database.db")
	like, err := db.Query("SELECT * FROM Likes WHERE user_id=" + strconv.Itoa(userInfo.ID))
	if err != nil {
		fmt.Println(err.Error())
	}

	var id int
	var idPost int
	var idUser int
	var dejaLike bool = false
	for like.Next() {
		err = like.Scan(&id, &idPost, &idUser)
		CheckErr(err)
		if idPost == upost_id {
			dejaLike = true
			break
		}
	}

	if dejaLike {
		likeNow = "checked"
	}

	like.Close()

	//Récupération du nouveau commentaire
	if r.Method == "POST" {
		db, _ := sql.Open("sqlite3", "database/database.db")
		comment := r.FormValue("comment")

		if userInfo.UserName != "" {
			if r.FormValue("Liker") != "" {
				if !dejaLike {
					datab, err := db.Prepare("INSERT INTO Likes (user_id,post_id) VALUES (?,?)")
					if err != nil {
						fmt.Println(err)
						http.Error(w, "Server Error", 500)
					}
					user_id := userInfo.ID
					post_id := upost_id
					_, err = datab.Exec(user_id, post_id)
					if err != nil {
						fmt.Println(err)
						fmt.Println("Une erreur ici")
					}
					datab.Close()
					fmt.Println("Tu viens de like ce post")

					dataPost, err := db.Query("SELECT * FROM Posts WHERE id=" + strconv.Itoa(post_id))
					if err != nil {
						fmt.Println(err.Error())
					}
					likes = 0
					var title string
					var categories string
					var body string
					var image string
					var comments_nb int
					for dataPost.Next() {
						err = dataPost.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb)
						CheckErr(err)
					}
					dataPost.Close()

					upost_id := strconv.Itoa(post_id)
					dataLikes, _ := db.Query("SELECT * FROM Likes WHERE post_id=" + upost_id)
					for dataLikes.Next() {
						err = dataLikes.Scan(&id, &upost_id, &user_id)
						CheckErr(err)
						likes++
					}
					dataLikes.Close()

					datab, err = db.Prepare("UPDATE Posts SET likes=? WHERE id=" + strconv.Itoa(post_id))
					fmt.Print("208 : ")
					fmt.Println(likes)
					datab.Exec(likes)
					datab.Close()
				} else {
					fmt.Println(" Tu as déjà like ce post")
				}
				likeNow = "checked"

			} else if comment != "" {
				datab, err := db.Prepare("INSERT INTO Comments (body, user_id,post_id) VALUES (?,?,?)")
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Server Error", 500)
				}

				fmt.Println(comment)
				user_id := userInfo.ID
				fmt.Println(user_id)
				post_id := upost_id
				fmt.Println(upost_id)
				fmt.Println("ICI")
				result, err := datab.Exec(comment, user_id, post_id)
				fmt.Println(result)
				if err != nil {
					fmt.Println(err)
					fmt.Println("Une erreur ici")
				}
				datab.Close()
			}
		} else {
			fmt.Println("LA")
			http.Redirect(w, r, "/login", 301)
		}
		db.Close()
	}

	db, _ = sql.Open("sqlite3", "database/database.db")
	//récupération de tout les commentaires liés au post
	fmt.Println("récupération de tout les commentaires liés au post")

	var title string
	var body string
	var image string
	var categories string
	var comments_nb int
	var allComments []COMMENT
	var user_id int

	getComment, err := db.Query("SELECT * FROM Comments WHERE post_id=" + post_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	var id_Post int
	var comment_id int
	var bodyComment string
	for getComment.Next() {
		err = getComment.Scan(&comment_id, &bodyComment, &user_id, &id_Post)
		CheckErr(err)
		if upost_id == id_Post {
			oneComment := COMMENT{
				ID:      comment_id,
				User_ID: user_id,
				Post_ID: id_Post,
				Body:    bodyComment,
			}

			userComment, err := db.Query("SELECT * FROM Users WHERE id=" + strconv.Itoa(user_id))
			if err != nil {
				fmt.Println(err.Error())
			}

			var user_comment INFO
			var id int
			var email string
			var password string
			var username string
			var since string
			var description string
			var image2 string
			var country string
			for userComment.Next() {
				err = userComment.Scan(&id, &username, &email, &since, &description, &password, &image2, &country)
				CheckErr(err)
				if id == user_id {
					user_comment = INFO{
						ID:          id,
						Email:       email,
						PassWord:    password,
						UserName:    username,
						Since:       since,
						Description: description,
						Image:       image2,
						Country:     country,
					}
					break
				}
			}
			userComment.Close()
			oneComment = COMMENT{
				ID:        comment_id,
				User_ID:   user_id,
				User_Info: user_comment,
				Post_ID:   id_Post,
				Body:      bodyComment,
			}

			allComments = append(allComments, oneComment)

		}
	}
	getComment.Close()

	//récupération du post
	fmt.Println("récupération du post")
	test, err := db.Query("SELECT * FROM Posts WHERE id=" + post_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	for test.Next() {
		err = test.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb)
		CheckErr(err)
		fmt.Print("test.Next: ")
		fmt.Println(likes)
	}
	test.Close()

	fmt.Println(title)
	tabCategories := strings.Split(categories, ";")
	var tabCat []CATEGORIES
	for _, x := range tabCategories {
		oneCategorie := CATEGORIES{
			Cat: x,
		}
		tabCat = append(tabCat, oneCategorie)
	}

	//Recupération des user_info du user qui a posté
	fmt.Println("Recupération des user_info du user qui a posté")
	user, err := db.Query("SELECT * FROM Users WHERE id=" + strconv.Itoa(user_id))
	if err != nil {
		fmt.Println(err.Error())
	}

	var post_user_info INFO
	var email string
	var password string
	var username string
	var since string
	var description string
	var image2 string
	var country string
	for user.Next() {
		err = user.Scan(&id, &username, &email, &since, &description, &password, &image2, &country)
		CheckErr(err)
		if id == user_id {
			post_user_info = INFO{
				ID:          id,
				Email:       email,
				PassWord:    password,
				UserName:    username,
				Since:       since,
				Description: description,
				Image:       image2,
				Country:     country,
			}
			break
		}
	}
	user.Close()

	post_info := POSTINFO{
		ID:           upost_id,
		User_ID:      post_user_info.ID,
		Title:        title,
		Body:         body,
		Image:        image,
		Categories:   tabCat,
		Likes:        likes,
		Comment_Nb:   comments_nb,
		All_Comments: allComments,
	}

	fmt.Println(post_info)

	data := ALLINFO{
		User_Info:           userInfo,
		Post_Info:           post_info,
		Post_User_Info:      post_user_info,
		Currently_Post_Like: likeNow,
	}

	fmt.Print("Nb Likes ")
	fmt.Println(data.Post_Info.Likes)

	defer db.Close()

	var files []string
	if data.Post_Info.Title != "" {
		files = []string{"template/Post.html", "template/Common.html"}
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

// AllPosts : kzdnzndz
func AllPosts(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

	allCategories := "sport;anime/manga;jeux vidéos;informatique;economie;voyage;NEWS;paranormal"
	tabCategories := strings.Split(allCategories, ";")
	var tabCat []CATEGORIES
	for _, x := range tabCategories {
		oneCategorie := CATEGORIES{
			Cat: x,
		}
		tabCat = append(tabCat, oneCategorie)
	}

	post := POSTINFO{
		AllCategories: tabCat,
	}

	data := ALLINFO{
		User_Info: user,
		Post_Info: post,
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
