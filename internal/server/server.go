package server

import (
	"bytes"
	"encoding/json"
	"github.com/VladKornilov/fondy_payments_go/internal/entities"
	"github.com/VladKornilov/fondy_payments_go/internal/fondy"
	"github.com/VladKornilov/fondy_payments_go/internal/logger"
	"github.com/google/uuid"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var app *Application
var templates map[string] string


func (a Application)StartServer() {
	app = &a
	err := app.db.Create()
	if logger.LogErr(err) { return }
	addPageTemplates()
	addPageListeners()
}

func addPageTemplates() {
	index, err := ioutil.ReadFile("html/templates/index.html")
	if logger.LogErr(err) { return }
	buy, err := ioutil.ReadFile("html/templates/buy.html")
	if logger.LogErr(err) { return }
	purchaseSuccess, err := ioutil.ReadFile("html/templates/purchase_success.html")
	if logger.LogErr(err) { return }

	templates = make(map[string] string)

	templates["index"] = string(index)
	templates["buy"] = string(buy)
	templates["purchaseSuccess"] = string(purchaseSuccess)
}

func addPageListeners() {
	response, _ := os.LookupEnv("RESPONSE_URL")
	port, _ := os.LookupEnv("SITE_PORT")
	http.HandleFunc("/buy", handleBuyRequest)
	http.HandleFunc("/purchase/", handleFondyRedirect)
	http.HandleFunc(response, handleResponse)
	http.HandleFunc("/", startPage)
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))

	err := http.ListenAndServeTLS(port, "ssl/server.crt", "ssl/server.key", nil)
	if logger.LogErr(err) { return }
}

func startPage(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("uuid")

	var userId string
	if err != nil {
		userId = uuid.New().String()
		err = app.db.InsertUser(entities.User{Uuid: userId})

		idCookie = new(http.Cookie)
		idCookie.Name = "uuid"
		idCookie.Value = userId
		idCookie.Expires = time.Now().Add(30 * 24 * time.Hour)
		http.SetCookie(w, idCookie)
		if logger.LogErr(err) { return }
	} else {
		userId = idCookie.Value
	}

	tpl, err := template.New("index").Parse(templates["index"])
	if logger.LogErr(err) { return }

	user, err := app.db.GetUserByUUID(userId)
	if logger.LogErr(err) { return }

	err = tpl.Execute(w, user)
	if logger.LogErr(err) { return }
}

func handleBuyRequest(w http.ResponseWriter, r *http.Request) {

	funcMap := template.FuncMap{
		"calcPrice": func (price int) string {
			whole := price / 100
			cents := price % 100
			zeroStr := ""
			if cents < 10 { zeroStr = "0" }
			return strconv.Itoa(whole) + "," + zeroStr + strconv.Itoa(cents)
		},
	}

	tpl, err := template.New("buy").Funcs(funcMap).Parse(templates["buy"])
	if logger.LogErr(err) { return }

	products := app.db.GetProducts()

	// 2) Торговец у себя на сайте отображает кнопку оплаты,
	//    при нажатии на которую будет осуществлен редирект
	//    пользователя на страницу ввода платежных реквизитов
	err = tpl.Execute(w,
		struct {
			Products []entities.Product
		} {
		products,
		})
	if logger.LogErr(err) { return }
}


// 4) торговец формирует host-to-host запрос на URL
//    https://pay.fondy.eu/api/checkout/url/, передавая параметры методом HTTPS POST
func handleFondyRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { return }
	//idCookie, err := r.Cookie("uuid")
	//if LogErr(err) { return }
	//userId := idCookie.Value

	var prodName string
	urlStr := r.URL.String()
	for i := len(urlStr)-1; urlStr[i] != '/'; i-- {
		if urlStr[i-1] == '/' {
			prodName = urlStr[i:]
		}
	}
	product, err := app.db.GetProductById(prodName)
	if logger.LogErr(err) { return }

	idCookie, err := r.Cookie("uuid")
	if err != nil { return }
	userId := idCookie.Value

	request := fondy.MakeRequest(userId, product)
	jsonObj := struct {
		Request fondy.Request `json:"request"`
	}{ request }

	data, err := json.Marshal(jsonObj)
	if logger.LogErr(err) { return }

	apiUrl, exists := os.LookupEnv("API_URL")
	if !exists { return }
	resp, err := http.Post(apiUrl, "application/json", bytes.NewReader(data))
	if logger.LogErr(err) { return }
	defer resp.Body.Close()

	// 5) Платежный шлюз FONDY возвращает торговцу промежуточный ответ
	//    с параметром checkout_url, содержащий URL, на который нужно
	//    перенаправить покупателя для ввода платежных реквизитов
	interResp := fondy.IntermediateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if logger.LogErr(err) { return }
	err = json.Unmarshal(body, &interResp)
	if logger.LogErr(err) { return }

	if interResp.Response.ResponseStatus == "failure" {
		logger.LogData("Failure: " + interResp.Response.ErrorMessage)
		return
	}
	logger.LogData("Inter Response success: Checkout URL = " + interResp.Response.CheckoutUrl)

	// 6) Торговец перенаправляет покупателя на checkout_url
	checkoutUrl := interResp.Response.CheckoutUrl
	http.Redirect(w, r, checkoutUrl, http.StatusSeeOther)

	// Покупатель вводит платежные реквизиты на сайте платежного шлюза FONDY,
	// платежный шлюз осуществляет списание средств через внешнюю платежную систему (банк-эквайер)
}


func handleResponse(w http.ResponseWriter, r *http.Request) {
	//if r.Method != "POST" { return }

	body, err := ioutil.ReadAll(r.Body)
	if logger.LogErr(err) { return }
	logger.LogData(string(body))

	values, err := url.ParseQuery(string(body))

	response := fondy.GetFinalResponse(values)
	productId, _ := strconv.Atoi(response.ProductId)

	product, err := app.db.GetProductById(response.ProductId)
	if logger.LogErr(err) { return }

	userId := response.MerchantData
	user, err := app.db.GetUserByUUID(userId)
	if logger.LogErr(err) { return }
	user.Diamonds += product.Value
	err = app.db.UpdateUser(user)
	if logger.LogErr(err) { return }

	purchase := entities.Purchase{
		PurchaseId:     response.OrderId,
		ProductId:      productId,
		UserId:			user.UserId,
		PurchaseStatus: response.ResponseStatus,
	}
	err = app.db.InsertPurchase(purchase)
	if logger.LogErr(err) { return }

	// 12) Торговец у себя на сайте отображает страницу с результатом оплаты

	tpl, err := template.New("purchaseSuccess").Parse(templates["purchaseSuccess"])
	if logger.LogErr(err) { return }
	
	err = tpl.Execute(w,
		struct {
		Amount int
		} {
			product.Value,
		})
	if logger.LogErr(err) { return }
}