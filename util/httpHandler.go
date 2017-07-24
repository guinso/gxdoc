package util

import (
	"fmt"
	"net/http"
	"strings"
)

// Generate error page
func HandleErrorCode(errorCode int, description string, w http.ResponseWriter) {
	w.WriteHeader(errorCode)                    // set HTTP status code (example 404, 500)
	w.Header().Set("Content-Type", "text/html") // clarify return type (MIME)
	w.Write([]byte(fmt.Sprintf(
		"<html><body><h1>Error %d</h1><p>%s</p></body></html>",
		errorCode,
		description)))
}

func HandleResponse(statusCode int, statusMsg string, jsonString string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json") //mime to application/json
	w.WriteHeader(http.StatusOK)                       //status code 200, OK
	w.Write([]byte(fmt.Sprintf(
		"{\"statusCode\" : %d, \"statusMsg\" : \"%s\", \"body\" : \"%s\"}",
		statusCode, statusMsg, jsonString))) //body text
}

func AddBackSlash(jsonString string) string {
	return strings.Replace(jsonString, "\"", "\\\"", -1)
}
