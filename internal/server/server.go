package server

import (
	"github.com/VladKornilov/fondy_payments_go/internal/entities"
	"github.com/google/uuid"
	"html/template"
	"io/ioutil"
	"net/http"
)

var app *Application
func (a Application)StartServer() {
	app = &a
	app.db.InsertProduct(entities.Product{ProductName: "test2", Price: 2000})
	addPageListeners()
}

func addPageListeners() {
	http.HandleFunc("/", startPage)
	http.HandleFunc("/buy", buyPage)
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))

	err := http.ListenAndServe(":8888", nil)
	if logErr(err) { return }
}

func startPage(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("uuid")

	var id string
	if err != nil {
		id = uuid.New().String()
		err = app.db.InsertUser(entities.User{Uuid: id})
		if logErr(err) { return }
	} else {
		id = idCookie.Value
	}

	bytes, err := ioutil.ReadFile("html/templates/index.html")
	if logErr(err) { return }

	tpl, err := template.New("index").Parse(string(bytes))
	if logErr(err) { return }

	user, err := app.db.GetUserByUUID(id)
	if logErr(err) { return }

	err = tpl.Execute(w, user)
	if logErr(err) { return }
}

func buyPage(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadFile("html/templates/buy.html")
	if logErr(err) { return }

	tpl, err := template.New("buy").Parse(string(bytes))
	if logErr(err) { return }

	products := app.db.GetProducts()

	data := struct {
		Products []entities.Product
	} {
		products,
	}

	err = tpl.Execute(w, data)
	if logErr(err) { return }
}