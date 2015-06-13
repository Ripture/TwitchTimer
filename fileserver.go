package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"path"
)

type Streamer struct {
	Name string
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Twitch)

	r.HandleFunc("/next", Next)

	fmt.Println("Starting server on :1935")
	http.ListenAndServe(":1935", r)
}

func Twitch(w http.ResponseWriter, r *http.Request) {
	ExeTemplate("nL_kripp", w, r)
}

func Next(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Taking Next")
	r.ParseForm()
	name := r.PostFormValue("streamer")
	fmt.Println(name)
	ExeTemplate(name, w, r)
}

func ExeTemplate(name string, w http.ResponseWriter, r *http.Request) {
	stre := Streamer{name}

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, stre); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
