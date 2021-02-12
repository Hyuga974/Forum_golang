package main


func main() {
	fs := http.FileServer(http.Dir("./template/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/Posts", AllPosts)
	http.HandleFunc("/Post", OnePost)
	http.HandleFunc("/Profil", Profil)
	http.HandleFunc("/Login", Login)
	http.HandleFunc("/Register", Register)

	fmt.Println("Start... ")
	http.ListenAndServe(":8080", nil)

}

// servHome : //* Page d'acceuil
func serveHome(w http.ResponseWriter, r *http.Request) {
	files := []string{"template/Home.html", "template/Common.html"}

	tmp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error: Check template", 500)
	}

	err = tmp.Execute(w, )
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}

}

func AllPosts(w http.ResponseWriter, r *http.Request){
	files := []string{"template/AllPosts.html", "template/Common.html"}
}

func OnePost(w http.ResponseWriter, r *http.Request){
	files := []string{"template/OnePost.html", "template/Common.html"}
}

func Profil(w http.ResponseWriter, r *http.Request){
	files := []string{"template/Profil.html", "template/Common.html"}
}

func Login(w http.ResponseWriter, r *http.Request){
	files := []string{"template/Login.html", "template/Common.html"}
}

func Register(w http.ResponseWriter, r *http.Request){
	files := []string{"template/Register.html", "template/Common.html"}
}