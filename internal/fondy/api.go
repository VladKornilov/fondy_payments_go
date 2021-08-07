package fondy

import (
	"github.com/VladKornilov/fondy_payments_go/internal/entities"
	"github.com/google/uuid"
	"net/url"
	"os"
	"strconv"
)

var Currencies = []string{"UAH", "RUB", "USD", "EUR", "GBP", "CZK"}
var VerificationStatuses = []string{"verified", "incorrect", "failed", "created"}

// Request formatting to json:
// bytes, err := json.Marshal(request)
type Request struct {
	OrderId				string `json:"order_id"`
	MerchantId			int    `json:"merchant_id"`
	OrderDesc			string `json:"order_desc"`
	Signature			string `json:"signature"`
	Amount				int    `json:"amount"`
	Currency			string `json:"currency"`
	MerchantData		string `json:"merchant_data"`
	ResponseUrl			string `json:"response_url"`
	ServerCallbackUrl	string `json:"server_callback_url"`
	ProductId			string `json:"product_id"`
}
func MakeRequest(customerId string, product entities.Product) Request {
	uuid := uuid.New().String()
	merchantIdStr, _ := os.LookupEnv("MERCHANT_ID")
	merchantId, _ := strconv.Atoi(merchantIdStr)
	siteUrl, _ := os.LookupEnv("SITE_URL")
	port, _ := os.LookupEnv("SITE_PORT")
	responseUrl, _ := os.LookupEnv("RESPONSE_URL")
	callbackUrl, _ := os.LookupEnv("CALLBACK_URL")

	request := Request{
		OrderId:    uuid,
		MerchantId: merchantId,
		OrderDesc:  product.ProductName,
		Amount:     product.Price,
		Currency:   "USD",
		ResponseUrl: siteUrl + port + responseUrl,
		ServerCallbackUrl: siteUrl + port + callbackUrl,
		MerchantData: customerId,

		ProductId:  strconv.Itoa(product.ProductId),
	}
	request.Signature = CalculateSignature(request)
	return request
}


type FinalResponse struct {
	OrderId             string `json:"order_id"`
	MerchantId          int    `json:"merchant_id"`
	Amount              int    `json:"amount"`
	Currency            string `json:"currency"`
	OrderStatus         string `json:"order_status"`
	ResponseStatus      string `json:"response_status"`
	Signature           string `json:"signature"`
	TranType            string `json:"tran_type"`
	SenderCellPhone     string `json:"sender_cell_phone"`
	SenderAccount       string `json:"sender_account"`
	MaskedCard          string `json:"masked_card"`
	CardBin             int    `json:"card_bin"`
	CardType            string `json:"card_type"`
	RRN                 string `json:"rrn"`
	ApprovalCode        string `json:"approval_code"`
	ResponseCode        string `json:"response_code"`
	ResponseDescription string `json:"response_description"`
	ReversalAmount      int    `json:"reversal_amount"`
	SettlementAmount    int    `json:"settlement_amount"`
	SettlementCurrency  string `json:"settlement_currency"`
	OrderTime           string `json:"order_time"`
	SettlementDate      string `json:"settlement_date"`
	ECI                 int    `json:"eci"`
	Fee                 int    `json:"fee"`
	PaymentSystem       string `json:"payment_system"`
	SenderEmail         string `json:"sender_email"`
	PaymentId           int    `json:"payment_id"`
	ActualAmount        int    `json:"actual_amount"`
	ActualCurrency      string `json:"actual_currency"`
	ProductId           string `json:"product_id"`
	MerchantData        string `json:"merchant_data"`
	VerificationStatus  string `json:"verification_status"`
	Rectoken            string `json:"rectoken"`
	RectokenLifetime    string `json:"rectoken_lifetime"`

	// error response
	ErrorCode      int    `json:"error_code"`
	ErrorMessage   string `json:"error_message"`
	RequestId	   string `json:"request_id"`
}

func GetFinalResponse(values url.Values) FinalResponse {
	amount, _ := strconv.Atoi(values["amount"][0])
	merchantId, _ := strconv.Atoi(values["merchant_id"][0])
	cardBin, _ := strconv.Atoi(values["card_bin"][0])
	revAmount, _ := strconv.Atoi(values["reversal_amount"][0])
	setAmount, _ := strconv.Atoi(values["settlement_amount"][0])
	eci, _ := strconv.Atoi(values["eci"][0])
	fee, _ := strconv.Atoi(values["fee"][0])
	paymentId, _ := strconv.Atoi(values["payment_id"][0])
	actualAmount, _ := strconv.Atoi(values["actual_amount"][0])

	return FinalResponse{
		OrderId: 				values["order_id"][0],
		MerchantId:				merchantId,
		Amount: 				amount,
		Currency: 				values["currency"][0],
		OrderStatus: 			values["order_status"][0],
		ResponseStatus: 		values["response_status"][0],
		Signature: 				values["signature"][0],
		TranType: 				values["tran_type"][0],
		SenderCellPhone: 		values["sender_cell_phone"][0],
		SenderAccount: 			values["sender_account"][0],
		MaskedCard: 			values["masked_card"][0],
		CardBin: 				cardBin,
		CardType: 				values["card_type"][0],
		RRN: 					values["rrn"][0],
		ApprovalCode: 			values["approval_code"][0],
		ResponseCode: 			values["response_code"][0],
		ResponseDescription: 	values["response_description"][0],
		ReversalAmount: 		revAmount,
		SettlementAmount:		setAmount,
		SettlementCurrency: 	values["settlement_currency"][0],
		OrderTime: 				values["order_time"][0],
		SettlementDate: 		values["settlement_date"][0],
		ECI: 					eci,
		Fee: 					fee,
		PaymentSystem: 			values["payment_system"][0],
		SenderEmail: 			values["sender_email"][0],
		PaymentId: 				paymentId,
		ActualAmount: 			actualAmount,
		ActualCurrency: 		values["actual_currency"][0],
		ProductId: 				values["product_id"][0],
		MerchantData: 			values["merchant_data"][0],
		VerificationStatus: 	values["verification_status"][0],
		Rectoken: 				values["rectoken"][0],
		RectokenLifetime: 		values["rectoken_lifetime"][0],
	}
}

type IntermediateResponse struct {
	Response Response
}

type Response struct {
	ResponseStatus string `json:"response_status"`
	CheckoutUrl    string `json:"checkout_url"`
	PaymentId      string `json:"payment_id"`

	// error response
	ErrorCode      int    `json:"error_code"`
	ErrorMessage   string `json:"error_message"`
	RequestId	   string `json:"request_id"`
}