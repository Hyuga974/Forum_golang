package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

//POSTINFO: Informations pour un Post
type ALLINFO struct {
	User_Info      INFO
	Post_Info      POSTINFO
	Post_User_Info INFO
}

//INFO: Déstiné à fournir des informations du user
type INFO struct {
	ID          int
	Email       string
	PassWord    string
	UserName    string
	Since       string
	Description string
	Image       string
	Country     string
	Login       bool
	Msg         string
}

type POSTINFO struct {
	ID           int
	User_ID      int
	Title        string
	Body         string
	Image        string
	Categories   []CATEGORIES
	Likes        int
	Comment_Nb   int
	All_Comments []COMMENT
}

type COMMENT struct {
	ID        int
	User_ID   int
	User_Info INFO
	Post_ID   int
	Body      string
}

type CATEGORIES struct {
	Cat string
}

//COOKIE: cookie
type Cookie struct {
	Name    string
	Value   string
	Expires time.Time
}

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/posts", AllPosts)
	http.HandleFunc("/post", OnePost)
	http.HandleFunc("/newpost", CreationPost)
	http.HandleFunc("/profil", Profil)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/test", Test)

	fmt.Println("Start... ")
	http.ListenAndServe(":8080", nil)

}

func Test(w http.ResponseWriter, r *http.Request) {

	files := []string{"template/CreatePost.html", "template/Common.html"}

	user := getSession(r)

	var Post POSTINFO
	var tabCat []CATEGORIES

	//sport, anime/manga, economie, jeux vidéo, informatique, voyages, NEW, paranormal.
	allCategories := "sport;anime/manga;jeux vidéos;informatique;economie;voyage;NEWS;paranormal"
	tabCategories := strings.Split(allCategories, ";")

	db, err := sql.Open("sqlite3", "database/database.db")
	checkErr(err)

	if r.Method == "POST" {
		fmt.Println("ici")
		datab, err := db.Prepare("INSERT INTO Posts (title, categories, body, user_id, image) VALUES (?,?,?,?,?)")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server Error", 500)
		}
		title := r.FormValue("title")
		body := r.FormValue("body")
		image := r.FormValue("myFile")

		fmt.Printf("image= %T\n", image)

		var categoriesCheck string
		for _, categorie := range tabCategories {
			if r.FormValue(categorie) != "" {
				categoriesCheck += categorie + ";"
			}
		}

		if title != "" && body != "" && categoriesCheck != "" {
			user_id := user.ID
			_, err := datab.Exec(title, categoriesCheck, body, user_id, image)
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
		files = []string{"template/Post", "template/Common.html"}
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

// servHome : //* Page d'acceuil
func serveHome(w http.ResponseWriter, r *http.Request) {
	user := getSession(r)

	data := ALLINFO{
		User_Info: user,
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

func CreationPost(w http.ResponseWriter, r *http.Request) {
	user := getSession(r)

	var Post POSTINFO
	var tabCat []CATEGORIES

	//sport, anime/manga, economie, jeux vidéo, informatique, voyages, NEW, paranormal.
	allCategories := "sport;anime/manga;jeux vidéos;informatique;economie;voyage;NEWS;paranormal"
	tabCategories := strings.Split(allCategories, ";")

	db, err := sql.Open("sqlite3", "database/database.db")
	checkErr(err)

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

		test, err := db.Query("SELECT * FROM Posts WHERE user_id=" + uId + " ORDER BY id DESC LIMIT 1")
		fmt.Println(test)
		if err != nil {
			fmt.Println(err.Error())
		}
		var id string
		var categories string
		var user_id int
		for test.Next() {
			err = test.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb)
			checkErr(err)
		}
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
}

//OnePost : Page pour un seul post
func OnePost(w http.ResponseWriter, r *http.Request) {
	userInfo := getSession(r)

	post_id := r.FormValue("id")
	upost_id, err := strconv.Atoi(post_id)

	db, err := sql.Open("sqlite3", "database/database.db")

	//Récupération du nouveau commentaire
	if r.Method == "POST" {
		comment := r.FormValue("comment")
		if userInfo.UserName != "" {

			datab, err := db.Prepare("INSERT INTO Comments (body, user_id,post_id) VALUES (?,?,?)")
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Server Error", 500)
			}

			_, err = datab.Exec(comment, userInfo.ID, upost_id)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			http.Redirect(w, r, "/login", 301)
		}
	}

	//récupération de tout les commentaires liés au post
	fmt.Println("récupération de tout les commentaires liés au post")

	var title string
	var body string
	var image string
	var categories string
	var likes int
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
		checkErr(err)
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
				checkErr(err)
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

	//récupération du post
	fmt.Println("récupération du post")
	test, err := db.Query("SELECT * FROM Posts WHERE id=" + post_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	for test.Next() {
		err = test.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb)
		checkErr(err)
	}

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
	var id int
	var email string
	var password string
	var username string
	var since string
	var description string
	var image2 string
	var country string
	for user.Next() {
		err = user.Scan(&id, &username, &email, &since, &description, &password, &image2, &country)
		checkErr(err)
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
		User_Info:      userInfo,
		Post_Info:      post_info,
		Post_User_Info: post_user_info,
	}

	fmt.Println(data.Post_Info.Title)

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
	db.Close()
}

// AllPosts : kzdnzndz
func AllPosts(w http.ResponseWriter, r *http.Request) {
	user := getSession(r)

	data := ALLINFO{
		User_Info: user,
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

func Profil(w http.ResponseWriter, r *http.Request) {
	user := getSession(r)

	data := ALLINFO{
		User_Info: user,
	}

	files := []string{"template/Profil.html", "template/Common.html"}

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

//
//
//	db, _ := sql.Open("sqlite3", "../database/database.db")
//	func OnePost(w http.ResponseWriter, r *http.Request) {
//	files := []string{"template/OnePost.html", "template/Common.html"}
//}
//
//func Profil(w http.ResponseWriter, r *http.Request) {
//	files := []string{"template/Profil.html", "template/Common.html"}
//}

//Register :  Permet de se créer à un compte
func Register(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database/database.db")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}
	msg := " "
	if r.Method == "POST" {
		datab, err := db.Prepare("INSERT INTO Users (username, email, since, description, password, image, country) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server Error", 500)
		}
		username := r.FormValue("username")
		email := r.FormValue("mail")
		loc, _ := time.LoadLocation("Europe/Paris")
		pretime := time.Now().In(loc)
		since := pretime.String()[:19]
		description := "Pas de description..."
		password := r.FormValue("password")
		image := "https://i.imgur.com/pMtf7R9.png"
		country := r.FormValue("country")
		confirm := r.FormValue("psw-confirmation")

		if username != "" || email != "" || password != "" {
			if password != confirm {
				msg = "Les deux mots de passe ne sont pas identiques"
			} else {
				_, err := datab.Exec(username, email, since, description, password, image, country)
				if err != nil {
					if err.Error() == "UNIQUE constraint failed: Users.email" {
						msg = "Cet E-Mail est déjà utilisé par un autre utilisateur"
					} else if err.Error() == "UNIQUE constraint failed: Users.username" {
						msg = "Ce nom est déjà utilisé par un autre utilisateur"
					} else {
						fmt.Println(err.Error())
					}
				}
			}
		}
	}
	Info := INFO{
		Msg: msg,
	}

	data := ALLINFO{
		User_Info: Info,
	}

	files := []string{"template/Register.html", "template/Common.html"}
	tmp, err := template.ParseFiles(files...) //err ne sert à rien!
	err = tmp.Execute(w, data)
	checkErr(err)

	db.Close()
}

// Login : Permet de se connecter à un compte
func Login(w http.ResponseWriter, r *http.Request) {
	userinfo := INFO{
		ID:          0,
		Email:       "",
		PassWord:    "",
		UserName:    "",
		Since:       "",
		Description: "",
		Image:       "",
		Country:     "",
		Login:       false,
		Msg:         "Vous êtes déconnecté",
	}

	db, err := sql.Open("sqlite3", "database/database.db")
	checkErr(err)

	cExist, id := checkSession(r)

	if r.Method == "POST" {
		if cExist {
			Delete(w, r, id)
		} else {
			test, err := db.Query("SELECT * FROM Users")
			if err != nil {
				fmt.Println(err.Error())
			}
			mailfound := false
			var id int
			var email string
			var password string
			var username string
			var since string
			var description string
			var image string
			var country string
			for test.Next() {
				err = test.Scan(&id, &username, &email, &since, &description, &password, &image, &country)
				checkErr(err)
				fmt.Println(id)
				fmt.Println(email)
				fmt.Println(password)
				if email == r.FormValue("mail") {
					mailfound = true
					break
				}
			}
			test.Close()
			if mailfound {
				if password == r.FormValue("password") {
					cookieCreation(w, id)
					userinfo = INFO{
						ID:          id,
						Email:       email,
						PassWord:    password,
						UserName:    username,
						Since:       since,
						Description: description,
						Image:       image,
						Country:     country,
						Login:       true,
						Msg:         "Vous êtes connecté en tant que " + username,
					}

				} else {
					userinfo = INFO{
						Msg: "Le mot de passe est invalide",
					}
				}
			} else {
				userinfo = INFO{
					Msg: "Ce mail n'est pas enregistré: veuillez vous inscrire",
				}
			}
		}

	} else {
		if cExist {
			userinfo = getSession(r)
		}
	}

	data := ALLINFO{
		User_Info: userinfo,
	}
	files := []string{"template/Connexion.html", "template/Common.html"}
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
	db.Close()
}

// Delete : "Supprime le UUID du compte qui se déconnect"
func Delete(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	checkErr(err)

	stmt, err := db.Prepare("delete from sessions where user_id=?")
	checkErr(err)

	res, err := stmt.Exec(id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

	c := http.Cookie{Name: "sessionLog", Value: "", MaxAge: -1}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/login", 301)
}

func cookieCreation(w http.ResponseWriter, id int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	checkErr(err)
	u1 := uuid.NewV4()
	fmt.Printf("UUID : %s ; User_id : %d\n", u1, id)
	datab, err := db.Prepare("INSERT INTO sessions (user_id, uuid) VALUES (?, ?)")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}
	_, err = datab.Exec(id, u1)
	if err != nil {
		fmt.Println(err)
	}
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "sessionLog", Value: String(u1), Expires: expiration}
	http.SetCookie(w, &cookie)
}

func checkSession(r *http.Request) (bool, int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	checkErr(err)
	cookie, _ := r.Cookie("sessionLog")
	ok := false
	var id int
	var uuid string
	if cookie != nil {
		dataCookie, err := db.Query("SELECT * FROM sessions")
		checkErr(err)

		for dataCookie.Next() {
			err = dataCookie.Scan(&id, &uuid)
			checkErr(err)
			fmt.Printf("id: %d\n", id)
			if cookie.Value == uuid {
				ok = true
				break
			}
		}
		dataCookie.Close()
	}
	fmt.Println(ok)
	return ok, id
}

func getSession(r *http.Request) INFO {
	var userinfo INFO
	cExist, idSession := checkSession(r)
	if cExist {
		db, err := sql.Open("sqlite3", "database/database.db")
		checkErr(err)
		tabusers, err := db.Query("SELECT * FROM Users")
		if err != nil {
			fmt.Println(err.Error())
		}
		var id int
		var email string
		var password string
		var username string
		var since string
		var description string
		var image string
		var country string
		for tabusers.Next() {
			err = tabusers.Scan(&id, &username, &email, &since, &description, &password, &image, &country)
			checkErr(err)
			if id == idSession {
				userinfo = INFO{
					ID:          id,
					Email:       email,
					PassWord:    password,
					UserName:    username,
					Since:       since,
					Description: description,
					Image:       image,
					Country:     country,
					Login:       true,
				}
				break
			}
		}
		tabusers.Close()
	}
	return userinfo
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("() ====================== ( ! ERROR ! ) ====================== ()")
		panic(err)
	}
}

func String(u uuid.UUID) string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}
