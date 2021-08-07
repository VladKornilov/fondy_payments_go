package main

import (
	"github.com/VladKornilov/fondy_payments_go/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	app, err := server.CreateApplication()
	
	if err != nil {
		panic(err)
	}
	defer func(app *server.Application) {
		err := app.Close()
		if err != nil {
			panic(err)
		}
	}(app)

	app.StartServer()
}
