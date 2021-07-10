package server

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

func (app Application)StartServer() {
	addPageListeners()
}

func addPageListeners() {
	http.HandleFunc("/", startPage)
	http.HandleFunc("/buy", buyPage)
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		println(err.Error())
		return
	}
}

func startPage(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("html/templates/index.html")
	if err != nil {
		println(err.Error())
		return
	}
	tpl, err := template.New("index").Parse(string(b))
	if err != nil {
		println(err.Error())
		return
	}

	product := struct {
		Title string
	} {
		Title: "Paradise Grind",
	}

	err = tpl.Execute(w, product)
	if err != nil {
		println(err.Error())
		return
	}
}

func buyPage(w http.ResponseWriter, r *http.Request) {

}