package fondy

import (
	"crypto/sha1"
	"fmt"
	"github.com/VladKornilov/fondy_payments_go/internal/logger"
	"os"
	"reflect"
	"sort"
)

// CalculateSignature creates a signature based on parameters given by Request
// signature = merch_pass | amount | currency | merch_id | order_desc | order_id
func CalculateSignature(request Request) string {
	merchPass, exists := os.LookupEnv("MERCHANT_PASSWORD")
	if !exists { return "" }

	t := reflect.TypeOf(request)
	names := make([]string, t.NumField())
	for i := range names {
		names[i] = t.Field(i).Name
	}
	sort.Strings(names)

	v := reflect.ValueOf(request)

	sign := merchPass
	for _, name := range names {
		val := v.FieldByName(name)
		fieldValue := fmt.Sprintf("%v", val.Interface())
		if fieldValue != "" {
			sign += "|" + fieldValue
		}
	}

	logger.LogData("Signature: " + sign)
	hash := sha1.New()
	hash.Write([]byte(sign))
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum)
}
