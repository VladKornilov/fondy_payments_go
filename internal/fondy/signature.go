package fondy

import (
	"crypto/sha1"
	"os"
	"strconv"
)

// CalculateSignature creates a signature based on parameters given by Request
// signature = merch_pass | amount | currency | merch_id | order_desc | order_id
func CalculateSignature(amount int, currency string, orderDesc string, orderId string) []byte {
	merchPass, exists := os.LookupEnv("MERCHANT_PASSWORD")
	if !exists { return nil }
	merchId, exists := os.LookupEnv("MERCHANT_ID")
	if !exists { return nil }

	sign := merchPass + "|" + strconv.Itoa(amount) + "|" +
		currency + "|" + merchId + "|" + orderDesc + "|" + orderId

	println("sign: " + sign)
	hash := sha1.New()
	hash.Write([]byte(sign))
	sum := hash.Sum(nil)
	return sum
}
