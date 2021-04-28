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

	_, err = res.RowsAffected()
	CheckErr(err)

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
		var mod int

		for tabusers.Next() {
			err = tabusers.Scan(&id, &username, &email, &since, &description, &password, &image, &country, &mod)
			CheckErr(err)

			if id == idSession {

				user := GetUser(id)
				posts := GetPost(user)
				userinfo = INFO{
					ID:          id,
					Email:       user.Email,
					PassWord:    user.PassWord,
					UserName:    user.UserName,
					Since:       user.Since,
					Description: user.Description,
					Image:       user.Image,
					Country:     user.Country,
					Admin:       user.Admin,
					Modo:        user.Modo,
					Login:       true,
					AllPosts:    posts,
				}
				break
			}
		}
		tabusers.Close()
		db.Close()
	}
	return userinfo
}

func GetUser(id int) INFO {
	fmt.Println("Récupération des info du user ", strconv.Itoa(id))
	db, err := sql.Open("sqlite3", "database/database.db")
	if err != nil {
		fmt.Print(err)
	}
	tabusers, err := db.Query("SELECT * FROM Users where id=" + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
	}
	var userinfo INFO
	var userAllPost []POSTINFO
	var userID int
	var username string
	var image string
	var email string
	var description string
	var password string
	var country string
	var since string
	var mod int
	for tabusers.Next() {
		err = tabusers.Scan(&userID, &username, &email, &since, &description, &password, &image, &country, &mod)
		CheckErr(err)
		if userID == id {
			userinfo = INFO{
				ID:          id,
				Email:       email,
				PassWord:    password,
				UserName:    username,
				Since:       since,
				Description: description,
				Image:       image,
				Country:     country,
				Mod:         mod,
			}
			userAllPost = GetPost(userinfo)

			break
		}
	}

	admin := IntToBoolAdmin(userinfo.Mod)
	modo := IntToBoolModo(userinfo.Mod)

	userinfo = INFO{
		ID:          id,
		Email:       email,
		PassWord:    password,
		UserName:    username,
		Since:       since,
		Description: description,
		Image:       image,
		Country:     country,
		Admin:       admin,
		Modo:        modo,
		AllPosts:    userAllPost,
	}
	tabusers.Close()
	db.Close()
	return userinfo
}

func IntToBoolAdmin(mod int) bool {
	if mod == 2 {
		return true
	} else {
		return false
	}
}

func IntToBoolModo(mod int) bool {
	if mod == 1 {
		return true
	} else {
		return false
	}
}

func GetPost(user INFO) []POSTINFO {
	var all_Post []POSTINFO
	db, err := sql.Open("sqlite3", "database/database.db")

	post, err := db.Query("SELECT * FROM Posts WHERE user_id=" + strconv.Itoa(user.ID) + " ORDER BY id DESC")

	color := RandomColor()
	var since string
	var user_id int
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
			ID:             idInt,
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
		"anime/manga":  "#D50C2E",
		"autre":        "#ccd1d1",
		"culture":      "#E15256",
		"economie":     "#C3C020",
		"informatique": "#19A9D1",
		"jeux vidéos":  "#23C009",
		"NEWS":         "#ff5733",
		"paranormal":   "#070709",
		"sport":        "#9D84C9",
		"voyage":       "#00FF12",
	}
	return allColor
}


func SearchData(SearchCriteria string) []int {
	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)

	var allUsers []INFO

	users, err := db.Query("SELECT * FROM Users ORDER BY id DESC")

	var currentlyUser INFO
	var email string
	var password string
	var username string
	var description string
	var country string
	var mod int
	var id int
	var image string
	var since string
	for users.Next() {
		err = users.Scan(&id, &username, &email, &since, &description, &password, &image, &country, &mod)
		CheckErr(err)
		currentlyUser = GetUser(id)
		allUsers = append(allUsers, currentlyUser)
	}
	users.Close()


	posts, _ := db.Query("SELECT * FROM Posts WHERE title LIKE '%"+SearchCriteria +"%'")

	color := RandomColor()
	var all_Post []POSTINFO
	var user_id int
	var title string
	var body string
	var likes int
	var comment_nb int
	var categories string
	var userinfo INFO
	for posts.Next() {
		err = posts.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
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
			userinfo = GetUser(user_id)

			post_info := POSTINFO{
				ID:         id,
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
	posts.Close()



	db.Close()



    var res []int
    // SearchCriteria = strings.ToLower(SearchCriteria)

    // for i, artist := range artistData {
    //     if strings.Contains(strings.ToLower(artist.Name), SearchCriteria) {
    //         res = append(res, i)
    //         continue
    //     }

    // }
    return res
}