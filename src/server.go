package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

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
	//http.HandleFunc("/profil", Profil)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/test", Test)

	fmt.Println("Start... ")
	http.ListenAndServe(":8080", nil)

}

func Test(w http.ResponseWriter, r *http.Request) {
	files := []string{"template/test.html", "template/Common.html"}

	tmp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error: Check template", 500)
	}

	err = tmp.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}

}

// servHome : //* Page d'acceuil
func serveHome(w http.ResponseWriter, r *http.Request) {
	user := getSession(r)

	files := []string{"template/Home.html", "template/Common.html"}

	tmp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error: Check template", 500)
	}

	err = tmp.Execute(w, user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}

}

//OnePost : Page pour un seul post
func OnePost(w http.ResponseWriter, r *http.Request) {
	user := getSession(r)

	files := []string{"template/Post.html", "template/Common.html"}

	tmp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error: Check template", 500)
	}

	err = tmp.Execute(w, user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}
}

// AllPosts : kzdnzndz
func AllPosts(w http.ResponseWriter, r *http.Request) {
	user := getSession(r)

	files := []string{"template/Posts.html", "template/Common.html"}

	tmp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error: Check template", 500)
	}

	err = tmp.Execute(w, user)
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
	data := INFO{
		Msg: msg,
	}
	files := []string{"template/Register.html", "template/Common.html"}
	tmp, err := template.ParseFiles(files...) //err ne sert à rien!
	err = tmp.Execute(w, data)
	checkErr(err)
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
	files := []string{"template/Connexion.html", "template/Common.html"}
	tmp, err := template.ParseFiles(files...)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error: Check template", 500)
	}

	err = tmp.Execute(w, userinfo)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}
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
