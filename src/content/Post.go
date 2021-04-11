package content

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func CreationPost(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

	if user.UserName != "" {
		r.ParseMultipartForm(10 << 20) //max size 10Mb (5mb for the pf)
		var Post POSTINFO
		var tabCat []CATEGORIES

		color := RandomColor()

		//sport, anime/manga, economie, jeux vidéo, informatique, voyages, NEW, paranormal.
		allCategories := "sport;anime/manga;jeux vidéos;informatique;economie;voyage;NEWS;paranormal"
		tabCategories := strings.Split(allCategories, ";")

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)

		if r.Method == "POST" {
			fmt.Println("ici")
			datab, err := db.Prepare("INSERT INTO Posts (title, categories, body, user_id, image, likes, comment_nb, since) VALUES (?,?,?,?,?,?,?, ?)")
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Server Error", 500)
			}
			title := r.FormValue("title")
			body := r.FormValue("body")
			image := r.FormValue("myFile")
			file, handler, err := r.FormFile("myFile")
			if err != nil {
				fmt.Println("Error Retrieving the File")
				fmt.Println(err)
				return
			}
			defer file.Close()
			fmt.Printf("Uploaded File: %+v\n", strings.ReplaceAll(handler.Filename, " ", "-"))
			fmt.Printf("File Size: %+v\n", handler.Size)
			fmt.Printf("MIME Header: %+v\n", handler.Header)

			absPath, _ := filepath.Abs("../src/assets/posts/" + strings.ReplaceAll(handler.Filename, " ", "-"))

			resFile, err := os.Create(absPath)
			if err != nil {
				fmt.Print(w, err)
			}
			defer resFile.Close()

			io.Copy(resFile, file)
			defer resFile.Close()
			fmt.Print("File uploaded")

			likes := 0
			comment_nb := 0
			loc, _ := time.LoadLocation("Europe/Paris")
			pretime := time.Now().In(loc)
			since := pretime.String()[:19]
			var categoriesCheck string
			for _, categorie := range tabCategories {
				if r.FormValue(categorie) != "" {
					categoriesCheck += categorie + ";"
				}
			}

			if title != "" && body != "" && categoriesCheck != "" {
				user_id := user.ID

				_, err := datab.Exec(title, categoriesCheck, body, user_id, image, likes, comment_nb, since)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			tabCategoriesCheck := strings.Split(categoriesCheck, ";")
			for _, x := range tabCategoriesCheck {
				oneCategorie := CATEGORIES{
					Cat:   x,
					Color: color[x],
				}
				tabCat = append(tabCat, oneCategorie)
			}
			Post = POSTINFO{
				User_ID:    user.ID,
				Title:      title,
				Body:       body,
				Image:      image,
				Categories: tabCat,
				Since:      since,
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
				err = newPost.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
				CheckErr(err)
			}
			newPost.Close()
			fmt.Println("id post:  ", id)

			http.Redirect(w, r, "/post?id="+id, 301)
		} else {
			for _, x := range tabCategories {
				oneCategorie := CATEGORIES{
					Cat:   x,
					Color: color[x],
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
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	color := RandomColor()

	postID := r.FormValue("id")
	if user.UserName != "" {
		var post_id int
		var title string
		var categories string
		var body string
		var user_id int
		var image string
		var likes int
		var comments_nb int
		var since string

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)

		Post, err := db.Query("SELECT * FROM Posts WHERE id=" + postID)
		if err != nil {
			fmt.Println(err.Error())
		}
		for Post.Next() {
			err = Post.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
			CheckErr(err)
		}
		Post.Close()

		tabCategories := strings.Split(categories, ";")
		var tabCat []CATEGORIES
		for _, x := range tabCategories {
			if x != "" {
				oneCategorie := CATEGORIES{
					Cat:   x,
					Color: color[x],
				}
				tabCat = append(tabCat, oneCategorie)
			}
		}

		//sport, anime/manga, economie, jeux vidéo, informatique, voyages, NEW, paranormal.
		allCategories := "sport;anime/manga;jeux vidéos;informatique;economie;voyage;NEWS;paranormal"
		tabAllCategories := strings.Split(allCategories, ";")

		var tabAllCat []CATEGORIES
		for _, x := range tabAllCategories {
			var check string
			for _, y := range tabCategories {
				if y == x {
					check = "checked"
					break
				}
			}
			oneCategorie := CATEGORIES{
				Cat:   x,
				Color: color[x],
				Check: check,
			}
			tabAllCat = append(tabAllCat, oneCategorie)
		}

		post_info := POSTINFO{
			ID:            post_id,
			User_ID:       user_id,
			Title:         title,
			Body:          body,
			Image:         image,
			Categories:    tabCat,
			AllCategories: tabAllCat,
			Likes:         likes,
			Comment_Nb:    comments_nb,
		}
		fmt.Println(post_info.Categories)

		if user.ID == post_info.User_ID {
			var tabCat []CATEGORIES

			if r.Method == "POST" {
				fmt.Println("Modification d'un Post en cours")

				//Récupération des nouvelles entrés
				newTitle := r.FormValue("title")
				newBody := r.FormValue("body")
				newImage := r.FormValue("Image")
				var newCategories string
				for _, categorie := range tabAllCategories {
					if r.FormValue(categorie) != "" {
						newCategories += categorie + ";"
					}
				}
				tabCategoriesCheck := strings.Split(newCategories, ";")
				for _, x := range tabCategoriesCheck {
					oneCategorie := CATEGORIES{
						Cat:   x,
						Color: color[x],
					}
					tabCat = append(tabCat, oneCategorie)
				}
				var Title string
				if newTitle != "" {
					Title = newTitle
				} else {
					Title = title
				}
				var Body string
				if body != "" {
					Body = newBody
				} else {
					Body = body
				}
				var Categories string
				if newCategories != "" {
					Categories = newCategories
				} else {
					Categories = categories
				}
				var Image string
				if newImage != "" {
					Image = newImage
				} else {
					Image = image
				}

				edit, _ := db.Prepare("UPDATE Posts SET title=?, categories=?, body=?, image=? WHERE id=" + postID)

				_, err := edit.Exec(Title, Categories, Body, Image)
				if err != nil {
					fmt.Println(err.Error())
				}

				post_info = POSTINFO{
					ID:         post_id,
					User_ID:    user_id,
					Title:      Title,
					Body:       Body,
					Image:      Image,
					Categories: tabCat,
				}

				edit.Close()
				http.Redirect(w, r, "/post?id="+strconv.Itoa(post_info.ID), 301)
			}
			data := ALLINFO{
				User_Info: user,
				Post_Info: post_info,
			}

			files := []string{"template/EditPost.html", "template/Common.html"}
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
	} else {
		fmt.Println("LA")
		http.Redirect(w, r, "/login", 301)
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

	postID := r.FormValue("id")
	fmt.Println("Post id :" + postID)
	if user.UserName != "" {
		var post_id int
		var title string
		var categories string
		var body string
		var user_id int
		var image string
		var likes int
		var comments_nb int
		var since string

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)
		post, err := db.Query("SELECT * FROM Posts WHERE id=" + postID)
		if err != nil {
			fmt.Println(err.Error())
		}

		CheckErr(err)
		for post.Next() {
			err = post.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
			CheckErr(err)
		}
		post.Close()

		post_info := POSTINFO{
			ID:         post_id,
			User_ID:    user_id,
			Title:      title,
			Body:       body,
			Image:      image,
			Likes:      likes,
			Comment_Nb: comments_nb,
		}

		if user.ID == post_info.User_ID {
			fmt.Println("Supresion d'un Post en cours")

			del, _ := db.Prepare("DELETE from Posts WHERE id=?")

			res, err := del.Exec(strconv.Itoa(post_info.ID))
			CheckErr(err)

			affect, err := res.RowsAffected()
			CheckErr(err)

			fmt.Println(affect)

			del.Close()
			http.Redirect(w, r, "/posts", 301)
			data := ALLINFO{
				User_Info: user,
				Post_Info: post_info,
			}

			files := []string{"template/EditPost.html", "template/Common.html"}
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
	} else {
		fmt.Println("LA")
		http.Redirect(w, r, "/login", 301)
	}
}

//OnePost : Page pour un seul post
func OnePost(w http.ResponseWriter, r *http.Request) {
	userInfo := GetSession(r)
	color := RandomColor()

	likeNow := ""
	likes := 0
	commentnb := 0

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
	var since string
	var dejaLike bool = false
	for like.Next() {
		err = like.Scan(&id, &idPost, &idUser, &since)
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
			if comment != "" {
				datab, err := db.Prepare("INSERT INTO Comments (body, user_id,post_id,since) VALUES (?,?,?,?)")
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
				loc, _ := time.LoadLocation("Europe/Paris")
				pretime := time.Now().In(loc)
				since := pretime.String()[:19]
				result, err := datab.Exec(comment, user_id, post_id, since)
				fmt.Println(result)
				if err != nil {
					fmt.Println(err)
					fmt.Println("Une erreur ici")
				}
				datab.Close()
			} else {
				if r.FormValue("Suppr") != "" {
					fmt.Printf("Post id : %s", post_id)
					upost_id, err := strconv.Atoi(post_id)
					if userInfo.ID == upost_id {

						CheckErr(err)
						fmt.Printf("User id : %d", userInfo.ID)
						stmt, err := db.Prepare("delete from Posts where user_id=? AND post_id=?")
						CheckErr(err)

						res, err := stmt.Exec(userInfo.ID, upost_id)
						CheckErr(err)

						affect, err := res.RowsAffected()
						CheckErr(err)

						fmt.Println(affect)

						stmt.Close()
					} else {
						fmt.Println("Vous ne pouvez pas supprimer")
					}
				} else {
					if !dejaLike {
						fmt.Println("je suis dans le if")

						loc, _ := time.LoadLocation("Europe/Paris")
						pretime := time.Now().In(loc)
						since := pretime.String()[:19]
						datab, err := db.Prepare("INSERT INTO Likes (user_id,post_id,since) VALUES (?,?,?)")
						if err != nil {
							fmt.Println(err)
							http.Error(w, "Server Error", 500)
						}
						user_id := userInfo.ID
						post_id := upost_id
						_, err = datab.Exec(user_id, post_id, since)
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
							err = dataPost.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
							CheckErr(err)
						}
						dataPost.Close()

						likeNow = "checked"

					} else {
						fmt.Printf("Post id : %s", post_id)
						upost_id, err := strconv.Atoi(post_id)
						CheckErr(err)
						fmt.Printf("User id : %d", userInfo.ID)
						stmt, err := db.Prepare("delete from Likes where user_id=? AND post_id=?")
						CheckErr(err)

						res, err := stmt.Exec(userInfo.ID, upost_id)
						CheckErr(err)

						affect, err := res.RowsAffected()
						CheckErr(err)

						fmt.Println(affect)

						stmt.Close()

						fmt.Println("Tu viens de unlike ce post")

						likeNow = ""
					}
				}
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
	var deletable bool

	getComment, err := db.Query("SELECT * FROM Comments WHERE post_id=" + post_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	var id_Post int
	var comment_id int
	var bodyComment string
	for getComment.Next() {
		err = getComment.Scan(&comment_id, &bodyComment, &user_id, &id_Post, &since)
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

	likes = 0
	dataLikes, _ := db.Query("SELECT * FROM Likes WHERE post_id=" + post_id)
	for dataLikes.Next() {
		err = dataLikes.Scan(&id, &upost_id, &user_id, &since)
		CheckErr(err)
		likes++
		fmt.Printf("nb likes pendant la boucle: %d", likes)
	}
	dataLikes.Close()

	datab, err := db.Prepare("UPDATE Posts SET likes=? WHERE id=" + post_id)
	CheckErr(err)
	datab.Exec(likes)
	datab.Close()

	commentnb = 0
	dataComment, _ := db.Query("SELECT * FROM Comments WHERE post_id=" + post_id)
	for dataComment.Next() {
		err = dataComment.Scan(&id, &body, &user_id, &upost_id, &since)
		CheckErr(err)
		commentnb++
		fmt.Printf("nb commentaires pendant la boucle: %d", commentnb)
	}
	dataComment.Close()

	datab, err = db.Prepare("UPDATE Posts SET comment_nb=? WHERE id=" + post_id)
	CheckErr(err)
	datab.Exec(commentnb)
	datab.Close()

	//récupération du post
	fmt.Println("récupération du post")
	test, err := db.Query("SELECT * FROM Posts WHERE id=" + post_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	for test.Next() {
		err = test.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
		CheckErr(err)
	}
	test.Close()
	fmt.Println(body)

	tabCategories := strings.Split(categories, ";")
	var tabCat []CATEGORIES
	for _, x := range tabCategories {
		oneCategorie := CATEGORIES{
			Cat:   x,
			Color: color[x],
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

	if post_user_info.ID == userInfo.ID {
		deletable = true
	}
	post_info := POSTINFO{
		ID:             upost_id,
		User_ID:        post_user_info.ID,
		Title:          title,
		Body:           body,
		Image:          image,
		Categories:     tabCat,
		Likes:          likes,
		Comment_Nb:     comments_nb,
		All_Comments:   allComments,
		Post_User_Info: post_user_info,
		Deletable:      deletable,
	}

	data := ALLINFO{
		User_Info:           userInfo,
		Post_Info:           post_info,
		Currently_Post_Like: likeNow,
	}

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

	post, err := db.Query("SELECT * FROM Posts ORDER BY id DESC")

	var since string
	var id string
	var user_id int
	var title string
	var body string
	var image string
	var likes int
	var comment_nb int
	var categories string
	for post.Next() {
		err = post.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
		CheckErr(err)
		idInt, _ := strconv.Atoi(id)
		cat := strings.Split(categories, ";")
		var tabCategories []CATEGORIES
		for _, x := range cat {
			catephemere := CATEGORIES{
				Cat:   x,
				Color: color[x],
			}
			tabCategories = append(tabCategories, catephemere)
		}

		tabusers, err := db.Query("SELECT * FROM Users")
		if err != nil {
			fmt.Println(err.Error())
		}
		var userinfo INFO
		var userAllPost []POSTINFO
		var userID int
		var username string
		var email string
		var description string
		var password string
		var country string
		for tabusers.Next() {
			err = tabusers.Scan(&userID, &username, &email, &since, &description, &password, &image, &country)
			CheckErr(err)
			if userID == user_id {
				userAllPost = GetPost(userID)
				userinfo = INFO{
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

		post_info = POSTINFO{
			ID:         idInt,
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
	post.Close()

	var allUsers []INFO

	users, err := db.Query("SELECT * FROM Users ORDER BY id DESC")
	var currentlyUser INFO
	var email string
	var password string
	var username string
	var description string
	var country string

	for users.Next() {
		err = users.Scan(&id, &username, &email, &since, &description, &password, &image, &country)
		CheckErr(err)

		idInt, _ := strconv.Atoi(id)
		currentlyUser = INFO{
			ID:          idInt,
			Email:       email,
			PassWord:    password,
			UserName:    username,
			Since:       since,
			Description: description,
			Image:       image,
			Country:     country,
		}
		allUsers = append(allUsers, currentlyUser)
	}
	users.Close()

	db.Close()

	postInfo := POSTINFO{
		AllCategories: tabCat,
	}

	data := ALLINFO{
		User_Info: user,
		Post_Info: postInfo,

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
