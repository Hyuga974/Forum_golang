package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Profil(w http.ResponseWriter, r *http.Request) {

	user := GetSession(r)
	if r.Method == "POST" {
		old_Description := user.Description
		old_Username := user.UserName
		old_Country := user.Country
		var new_Username string
		var new_Description string
		var new_Country string
		if r.FormValue("Username") != "" {
			new_Username = r.FormValue("Username")
		} else {
			new_Username = old_Username
		}
		if r.FormValue("Description") != "" {
			new_Description = r.FormValue("Description")
		} else {
			new_Description = old_Description
		}

		if r.FormValue("country") != "" {
			new_Country = r.FormValue("country")
		} else {
			new_Country = old_Country
		}

		db, _ := sql.Open("sqlite3", "database/database.db")
		datab, err := db.Prepare("UPDATE Users SET username=?, description=?, country=? WHERE id=" + strconv.Itoa(user.ID))
		CheckErr(err)
		datab.Exec(new_Username, new_Description, new_Country)
		datab.Close()

		user = GetSession(r)
	}
	db, _ := sql.Open("sqlite3", "database/database.db")

	//Récupèrelast post
	lastpost, err := db.Query("SELECT * FROM Posts where user_id=" + strconv.Itoa(user.ID) + " ORDER BY id DESC LIMIT 1")
	if err != nil {
		fmt.Println(err.Error())
	}
	color := RandomColor()
	var post_info_last_posted POSTINFO
	var id int
	var user_id int
	var title string
	var body string
	var image string
	var likes int
	var comment_nb int
	var categories string
	var since string
	for lastpost.Next() {
		err = lastpost.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
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
		post_info_last_posted = POSTINFO{
			ID:             id,
			User_ID:        user_id,
			Title:          title,
			Body:           body,
			Image:          image,
			Categories:     tabCategories,
			Likes:          likes,
			Comment_Nb:     comment_nb,
			Since:          since,
			Post_User_Info: user,
		}
		fmt.Println(post_info_last_posted)
	}
	lastpost.Close()

	//Récupération Like
	var post_info_last_like POSTINFO
	lastlike, err := db.Query("SELECT * FROM Likes where user_id=" + strconv.Itoa(user.ID) + " ORDER BY id DESC LIMIT 1")
	if err != nil {
		fmt.Println(err.Error())
	}
	var post_id int
	for lastlike.Next() {
		err = lastlike.Scan(&id, &post_id, &user_id, &since)
		CheckErr(err)
	}
	lastpostlike, err := db.Query("SELECT * FROM Posts where id=" + strconv.Itoa(post_id))
	if err != nil {
		fmt.Println(err.Error())
	}
	for lastpostlike.Next() {
		err = lastpostlike.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
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
		post_info_last_like = POSTINFO{
			ID:             id,
			User_ID:        user_id,
			Title:          title,
			Body:           body,
			Image:          image,
			Categories:     tabCategories,
			Likes:          likes,
			Comment_Nb:     comment_nb,
			Since:          since,
			Post_User_Info: user,
		}
	}
	lastpostlike.Close()
	lastlike.Close()

	//Récupération Last post Comment
	var post_info_last_comment POSTINFO
	lastcomment, err := db.Query("SELECT * FROM Comments where user_id=" + strconv.Itoa(user.ID) + " ORDER BY id DESC LIMIT 1")
	if err != nil {
		fmt.Println(err.Error())
	}
	for lastcomment.Next() {
		err = lastcomment.Scan(&id, &body, &post_id, &user_id, &since)
		CheckErr(err)
	}
	lastcomment.Close()

	lastpostcomment, err := db.Query("SELECT * FROM Posts where id=" + strconv.Itoa(post_id))
	if err != nil {
		fmt.Println(err.Error())
	}
	for lastpostcomment.Next() {
		err = lastpostcomment.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
		CheckErr(err)

	}
	cat := strings.Split(categories, ";")
	var tabCategories []CATEGORIES
	for _, x := range cat {
		catephemere := CATEGORIES{
			Cat:   x,
			Color: color[x],
		}
		tabCategories = append(tabCategories, catephemere)
	}
	post_info_last_comment = POSTINFO{
		ID:             id,
		User_ID:        user_id,
		Title:          title,
		Body:           body,
		Image:          image,
		Categories:     tabCategories,
		Likes:          likes,
		Comment_Nb:     comment_nb,
		Since:          since,
		Post_User_Info: user,
	}
	fmt.Println("Last commentaire :", post_info_last_comment)
	lastpostcomment.Close()

	data := ALLINFO{
		User_Info:    user,
		Last_Post:    post_info_last_posted,
		Last_Like:    post_info_last_like,
		Last_Comment: post_info_last_comment,
	}

	fmt.Println(data.Last_Post)
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
	allcountry := getPays()
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
		Crypted := []byte(password)
		Crypted, _ = bcrypt.GenerateFromPassword(Crypted, 10)
		fmt.Printf("Username : %s ; Email: %s; Country: %s;", username, email, country)

		if username != "" || email != "" || password != "" {
			if password != confirm {
				msg = "Les deux mots de passe ne sont pas identiques"
			} else {
				_, err := datab.Exec(username, email, since, description, Crypted, image, country)
				if err != nil {
					if err.Error() == "UNIQUE constraint failed: Users.email" {
						msg = "Cet E-Mail est déjà utilisé par un autre utilisateur"
					} else if err.Error() == "UNIQUE constraint failed: Users.username" {
						msg = "Ce nom est déjà utilisé par un autre utilisateur"
					} else {
						fmt.Println(err.Error())
						http.Redirect(w, r, "/login", 301)
					}
				}
			}
			fmt.Println(msg)
			fmt.Println(" Compte créé ")
		}
	}
	Info := INFO{
		Msg: msg,
	}

	data := ALLINFO{
		User_Info:   Info,
		All_Country: allcountry,
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
			var Password string
			var username string
			var since string
			var description string
			var image string
			var country string
			for test.Next() {
				err = test.Scan(&id, &username, &email, &since, &description, &Password, &image, &country)
				CheckErr(err)
				fmt.Println(id)
				fmt.Println(email)
				fmt.Println(Password)
				fmt.Println(r.FormValue("mail"))
				if email == r.FormValue("mail") {
					mailfound = true
					break
				}
			}
			mdp := r.FormValue("password")
			fmt.Printf("mdp entré : %s", mdp)

			test.Close()
			fmt.Printf("Before: mailfound --> %v \n", mailfound)
			if mailfound {
				fmt.Print("Into mailfound")
				fmt.Println(mdp)
				fmt.Println(Password)
				cryptedPassword := []byte(Password)
				fmt.Println(bcrypt.CompareHashAndPassword(cryptedPassword, []byte(mdp)))
				if bcrypt.CompareHashAndPassword(cryptedPassword, []byte(mdp)) == nil {
					CookieCreation(w, id)
					userinfo = INFO{
						ID:          id,
						Email:       email,
						PassWord:    mdp,
						UserName:    username,
						Since:       since,
						Description: description,
						Image:       image,
						Country:     country,
						Login:       true,
						Msg:         "Vous êtes connecté en tant que " + username,
					}

					http.Redirect(w, r, "/", 301)
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
