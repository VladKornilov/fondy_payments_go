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
	"os"
	"strconv"
	"time"
)

var app *Application
func (a Application)StartServer() {
	app = &a
	app.db.InsertProduct(entities.Product{ProductName: "test2", Price: 2000})
	addPageListeners()
}

func addPageListeners() {

	http.HandleFunc("/buy", buyPage)
	http.HandleFunc("/purchase/", fondyRedirect)
	http.HandleFunc("/", startPage)
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))

	err := http.ListenAndServe(":8888", nil)
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

	bytes, err := ioutil.ReadFile("html/templates/index.html")
	if logger.LogErr(err) { return }

	tpl, err := template.New("index").Parse(string(bytes))
	if logger.LogErr(err) { return }

	user, err := app.db.GetUserByUUID(userId)
	if logger.LogErr(err) { return }

	err = tpl.Execute(w, user)
	if logger.LogErr(err) { return }
}

func buyPage(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("html/templates/buy.html")
	if logger.LogErr(err) { return }

	funcMap := template.FuncMap{
		"calcPrice": func (price int) string {
			whole := price / 100
			cents := price % 100
			return strconv.Itoa(whole) + "," + strconv.Itoa(cents)
		},
	}

	tpl, err := template.New("buy").Funcs(funcMap).Parse(string(data))
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
func fondyRedirect(w http.ResponseWriter, r *http.Request) {
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

	request := fondy.MakeRequest(app.Config.MerchantId, product)
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
}