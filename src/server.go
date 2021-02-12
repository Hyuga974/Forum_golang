package main


func main() {
	fs := http.FileServer(http.Dir("./template/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", serveHome)

	fmt.Println("Start... ")

	http.ListenAndServe(":8080", nil)

}

// servHome : //* Page d'acceuil
func serveHome(w http.ResponseWriter, r *http.Request) {
	files := []string{"template/Home.html", "template/Nav.html"}

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