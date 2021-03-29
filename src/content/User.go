package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

func Profil(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

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
	CheckErr(err)

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
	CheckErr(err)

	cExist, id := CheckSession(r)

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
				CheckErr(err)
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
					CookieCreation(w, id)
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
			userinfo = GetSession(r)
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
