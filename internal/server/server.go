package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func (app Application)StartServer() {
	addPageListeners()
}

func addPageListeners() {
	http.HandleFunc("/", startPage)

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		println(err.Error())
		return
	}
}

func startPage(w http.ResponseWriter, r *http.Request) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	b, err := ioutil.ReadFile("internal/templates/index.hbs")
	if err != nil {
		println(err.Error())
		return
	}
	_, err = fmt.Fprintf(w, string(b))
	if err != nil {
		println(err.Error())
		return
	}
}