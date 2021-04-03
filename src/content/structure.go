package content

import "time"

//POSTINFO: Informations pour un Post
type ALLINFO struct {
	User_Info INFO
	Post_Info POSTINFO

	All_User            []INFO
	All_Country         []string
	All_Posts           []POSTINFO
	Post_Most_Recent    []POSTINFO
	Post_Most_Likes     []POSTINFO
	Currently_Post_Like string
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

	AllPosts []POSTINFO
}

type POSTINFO struct {
	ID            int
	User_ID       int
	Title         string
	Body          string
	Image         string
	Categories    []CATEGORIES
	AllCategories []CATEGORIES
	Likes         int
	Comment_Nb    int
	All_Comments  []COMMENT
	Since         string
	Deletable     bool

	Post_User_Info INFO
}

type COMMENT struct {
	ID        int
	User_ID   int
	User_Info INFO
	Post_ID   int
	Body      string
}

type CATEGORIES struct {
	Cat   string
	Color string
}

//COOKIE: cookie
type Cookie struct {
	Name    string
	Value   string
	Expires time.Time
}
