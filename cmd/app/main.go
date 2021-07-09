package main

import (
	"fmt"
	"github.com/VladKornilov/fondy_payments_go/internal/fondy"
	"github.com/VladKornilov/fondy_payments_go/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	//p := entities.Purchase{
	//	ProductId: 228,
	//	PurchaseStatus: "consumed",
	//}
	//err = db.InsertPurchase(p)
	//if err != nil {
	//	println(err.Error())
	//}
	//ps := db.GetPurchases()

	request := fondy.Request{
		OrderId: "test123456",
		OrderDesc: "test order",
		Currency: "USD",
		Amount: 125,
	}
	sign := fondy.CalculateSignature(request.Amount, request.Currency, request.OrderDesc, request.OrderId)
	fmt.Printf("%x", sign)

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
