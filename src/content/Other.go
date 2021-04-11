package content

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Delete : "Supprime le UUID du compte qui se déconnect"
func Delete(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)

	stmt, err := db.Prepare("delete from sessions where user_id=?")
	CheckErr(err)

	res, err := stmt.Exec(id)
	CheckErr(err)

	affect, err := res.RowsAffected()
	CheckErr(err)

	fmt.Println(affect)

	db.Close()

	c := http.Cookie{Name: "sessionLog", Value: "", MaxAge: -1}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/login", 301)
}

func CookieCreation(w http.ResponseWriter, id int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)
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

	db.Close()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "sessionLog", Value: String(u1), Expires: expiration}
	http.SetCookie(w, &cookie)
}

func CheckSession(r *http.Request) (bool, int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)
	cookie, _ := r.Cookie("sessionLog")
	ok := false
	var id int
	var uuid string
	if cookie != nil {
		dataCookie, err := db.Query("SELECT * FROM sessions")
		CheckErr(err)

		for dataCookie.Next() {
			err = dataCookie.Scan(&id, &uuid)
			CheckErr(err)
			if cookie.Value == uuid {
				ok = true
				break
			}
		}
		dataCookie.Close()
	}

	db.Close()
	return ok, id
}

func GetSession(r *http.Request) INFO {
	var userinfo INFO
	color := RandomColor()
	fmt.Println("Get Session")
	cExist, idSession := CheckSession(r)
	if cExist {

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)

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

		var all_Post []POSTINFO

		for tabusers.Next() {
			err = tabusers.Scan(&id, &username, &email, &since, &description, &password, &image, &country)
			CheckErr(err)

			if id == idSession {

				var postID int
				var title string
				var categorie string
				var body string
				var userID int
				var postImage string
				var likes int
				var comment_nb int
				var since string

				post, _ := db.Query("SELECT * from Posts WHERE user_id=" + strconv.Itoa(id))
				for post.Next() {
					err = post.Scan(&postID, &title, &categorie, &body, &userID, &postImage, &likes, &comment_nb, &since)
					if id == userID {

						tabCategories := strings.Split(categorie, ";")
						var tabCat []CATEGORIES
						for _, x := range tabCategories {
							oneCategorie := CATEGORIES{
								Cat:   x,
								Color: color[x],
							}
							tabCat = append(tabCat, oneCategorie)
						}

						var post_user_info INFO
						var post_user_id int
						var post_user_email string
						var post_user_password string
						var post_user_username string
						var post_user_description string
						var post_user_since string
						var post_user_image2 string
						var post_user_country string
						user, _ := db.Query("SELECT * FROM Users Where id=" + strconv.Itoa(userID))
						for user.Next() {
							err = user.Scan(&post_user_id, &post_user_username, &post_user_email, &post_user_since, &post_user_description, &post_user_password, &post_user_image2, &post_user_country)
							CheckErr(err)
							if id == user_id {
								post_user_info = INFO{
									ID:          post_user_id,
									Email:       post_user_email,
									PassWord:    post_user_password,
									UserName:    poost_user_username,
									Since:       post_user_since,
									Description: post_user_description,
									Image:       post_user_image2,
									Country:     post_user_country,
								}
								break
							}
						}
						user.Close()
						post_info := POSTINFO{
							ID:         postID,
							User_ID:    id,
							Title:      title,
							Body:       body,
							Image:      postImage,
							Categories: tabCat,
							Likes:      likes,
							Comment_Nb: comment_nb,

							Post_User_Info: post_user_info,
						}

						all_Post = append(all_Post, post_info)
					}
				}
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
					AllPosts:    all_Post,
				}
				break
			}
		}
		tabusers.Close()
		db.Close()
	}
	fmt.Println("Get session finit")
	return userinfo
}

func GetPost(user_id int) []POSTINFO {
	var all_Post []POSTINFO
	db, err := sql.Open("sqlite3", "database/database.db")

	post, err := db.Query("SELECT * FROM Posts WHERE user_id=" + strconv.Itoa(user_id) + " ORDER BY id DESC")

	color := RandomColor()
	var since string
	var id string
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
		post_info := POSTINFO{
			ID:         idInt,
			User_ID:    user_id,
			Title:      title,
			Body:       body,
			Image:      image,
			Categories: tabCategories,
			Likes:      likes,
			Comment_Nb: comment_nb,
			Since:      since,
		}
		all_Post = append(all_Post, post_info)
	}
	post.Close()

	db.Close()
	return all_Post
}

func CheckErr(err error) {
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

func RandomColor() map[string]string {
	allColor := map[string]string{
		"informatique": "#19A9D1",
		"anime/manga":  "#D50C2E",
		"jeux vidéos":  "#23C009",
		"sport":        "#9D84C9",
		"economie":     "#C3C020",
		"voyage":       "#00FF12",
		"NEWS":         "#CDC8C6",
		"paranormal":   "#070709",
	}
	return allColor
}
