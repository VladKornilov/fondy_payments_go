package fondy

import (
	"github.com/VladKornilov/fondy_payments_go/internal/entities"
	"github.com/google/uuid"
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
	Version				string `json:"version"`
	MerchantData		string `json:"merchant_data"`
	ServerCallbackUrl	string `json:"server_callback_url"`
	ProductId			string `json:"product_id"`
}
func MakeRequest(merchantId int, product entities.Product) Request {
	uuid := uuid.New().String()
	request := Request{
		OrderId:    uuid,
		MerchantId: merchantId,
		OrderDesc:  product.ProductName,
		Amount:     product.Price,
		Currency:   "USD",
		ServerCallbackUrl: "https://127.0.0.1/purchase_server_callback_url",
		MerchantData: "our_custom_payload",

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
	AdditionalInfo      string `json:"additional_info"`

	// error response
	ErrorCode      int    `json:"error_code"`
	ErrorMessage   string `json:"error_message"`
	RequestId	   string `json:"request_id"`
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