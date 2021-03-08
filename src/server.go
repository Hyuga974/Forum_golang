package main

import (
	"database/sql"
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

type COOKIE struct {
	ID   int
	UUID int
}

func main() {
	fs := http.FileServer(http.Dir("./template/assets"))
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
	files := []string{"template/Home.html", "template/Common.html"}

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

//OnePost : Page pour un seul post
func OnePost(w http.ResponseWriter, r *http.Request) {
	files := []string{"template/Post.html", "template/Common.html"}

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

// AllPosts : kzdnzndz
func AllPosts(w http.ResponseWriter, r *http.Request) {

	files := []string{"template/Posts.html", "template/Common.html"}

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

// Login : Permet de se connecter à un compte
func Login(w http.ResponseWriter, r *http.Request) {

	msg := "Vous êtes déconnecté"

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
		Msg:         msg,
	}

	db, err := sql.Open("sqlite3", "database/database.db")
	checkErr(err)

	if r.Method == "POST" {
		// 	test, err := db.Query("SELECT " + r.FormValue("mail") + " FROM Users")
		// 	checkErr(err)
		// 	fmt.Println(test)

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
				msg = "Vous êtes connecté en tant que " + username
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
					Msg:         msg,
				}
				u1 := uuid.NewV4()
				fmt.Printf("UUID : %s ; User_id : %d\n", u1, userinfo.ID)
				datab, err := db.Prepare("INSERT INTO sessions (user_id, uuid) VALUES (?, ?)")
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Server Error", 500)
				}
				_, err = datab.Exec(userinfo.ID, u1)
				if err != nil {
					fmt.Println(err)
				}
				// cookie := http.COOKIE{ID: userinf.ID, Value: u1}
				// http.SetCookie(w, &cookie)
			} else {
				msg = "Le mot de passe est invalide"
				userinfo = INFO{
					Msg: msg,
				}
			}
		} else {
			msg = "Ce mail n'est pas enregistré: veuillez vous inscrire"
			userinfo = INFO{
				Msg: msg,
			}
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

// func checkSession() {

// }

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
