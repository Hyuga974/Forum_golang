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
)

func Profil(w http.ResponseWriter, r *http.Request) {

	user := GetSession(r)
	if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20) //max size 10Mb (5mb for the pf)
		old_Description := user.Description
		old_Username := user.UserName
		old_Country := user.Country
		old_Image := user.Image
		var new_Username string
		var new_Description string
		var new_Country string
		var new_Image string
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

		file, handler, err := r.FormFile("myFile")
		if err != nil {
			fmt.Println(err)
			new_Image = old_Image
		} else {
			defer file.Close()
			fmt.Printf("Uploaded File: %+v\n", strings.ReplaceAll(handler.Filename, " ", "-"))
			fmt.Printf("File Size: %+v\n", handler.Size)
			fmt.Printf("MIME Header: %+v\n", handler.Header)

			/*
				buff := make([]byte, 512)
				_, err = file.Read(buff)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				filetype := http.DetectContentType(buff)
				if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" {
					http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG or GIF image", http.StatusBadRequest)
					return
				}
			*/

			absPath, _ := filepath.Abs("../src/assets/profiles/" + strings.ReplaceAll(handler.Filename, " ", "-"))

			resFile, err := os.Create(absPath)
			if err != nil {
				fmt.Print(w, err)
			}
			defer resFile.Close()

			io.Copy(resFile, file)
			defer resFile.Close()
			fmt.Print("File uploaded")

			new_Image = "../assets/profiles/" + strings.ReplaceAll(handler.Filename, " ", "-")
		}

		db, _ := sql.Open("sqlite3", "database/database.db")
		datab, err := db.Prepare("UPDATE Users SET username=?, description=?, country=?, image=? WHERE id=" + strconv.Itoa(user.ID))
		CheckErr(err)
		datab.Exec(new_Username, new_Description, new_Country, new_Image)
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
		post_user := GetUser(user_id)
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
			Post_User_Info: post_user,
		}
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

		post_user := GetUser(user_id)
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
			Post_User_Info: post_user,
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
		err = lastcomment.Scan(&id, &body, &user_id, &post_id, &since)
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

	post_user := GetUser(user_id)
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
		Post_User_Info: post_user,
	}
	lastpostcomment.Close()

	allcountry := getPays()

	data := ALLINFO{
		All_Country:  allcountry,
		User_Info:    user,
		Last_Post:    post_info_last_posted,
		Last_Like:    post_info_last_like,
		Last_Comment: post_info_last_comment,
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
