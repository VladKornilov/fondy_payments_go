package server

import "time"

// logErr prints content of error to console
// with additional information.
// Returns true, if error exists
func logErr(err error) bool {
	if err != nil {
		message := time.RFC822 + " ERROR:\n" + err.Error()
		println(message)
		return true
	}
	return false
}

func logData(data string) {
	message := time.RFC822 + " DATA:\n" + data
	println(message)
}