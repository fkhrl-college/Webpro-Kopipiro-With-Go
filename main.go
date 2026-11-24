package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl = template.Must(
	template.New("").Funcs(template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}).ParseGlob("templates/*.html"))

func main() {
	connectDB()
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	mux.HandleFunc("GET /{$}", requireLogin(indexHandler))
	mux.HandleFunc("/create", requireLogin(createHandler))
	mux.HandleFunc("/edit", requireLogin(editHandler))
	mux.HandleFunc("/delete", requireLogin(deleteHandler))
	log.Println("Server running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
