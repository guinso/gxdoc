package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/guinso/stringtool"

	//explicitly include GO mysql library
	//_ "github.com/go-sql-driver/mysql"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)

var dbb *sql.DB

var productionDB *sql.DB

//SetDB set and hold production database handler
func SetDB(db *sql.DB) {
	productionDB = db
}

//GetDB get production database handler
func GetDB() *sql.DB {
	return productionDB
}

//DecodeJSON decode http request's JSON into golang's instance
func DecodeJSON(request *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(request.Body)

	return decoder.Decode(&obj)
}

//GetRandomRunningNumber get next random generated MD5 value to fill
func GetRandomRunningNumber(tableName string) string {

	return stringtool.MakeMD5(
		tableName +
			strconv.FormatInt(time.Now().UnixNano(), 10) +
			strconv.FormatInt(rand.Int63(), 10))
}

//SendHTTPResponse send HTTP response
func SendHTTPResponse(w http.ResponseWriter, statusCode int, statusMsg string, json string) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(
		"{\"statusCode\":%d, \"statusMsg\":\"%s\", \"response\":%s}",
		statusCode, statusMsg, json)))
}

//SendHTTPErrorResponse send HTTP 500 internal error response
func SendHTTPErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{msg:\"Encounter internal server error\"}"))
}

//IsPOST check request HTTP is POST method
func IsPOST(r *http.Request) bool {
	return strings.Compare(strings.ToLower(r.Method), "post") == 0
}

//IsGET check request HTTP is GET method
func IsGET(r *http.Request) bool {
	return strings.Compare(strings.ToLower(r.Method), "get") == 0
}

//IsPUT check request HTTP is PUT method
func IsPUT(r *http.Request) bool {
	return strings.Compare(strings.ToLower(r.Method), "put") == 0
}

//IsDELETE check request HTTP is DELETE method
func IsDELETE(r *http.Request) bool {
	return strings.Compare(strings.ToLower(r.Method), "delete") == 0
}

//IsFileExists check specify file path is exists or not
func IsFileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false //file not found
		}

		return false //stat command error
	}

	return true //file exists
}

//IsDirectoryExists check specify directory path is exists or not
func IsDirectoryExists(directoryName string) (bool, error) {
	stat, err := os.Stat(directoryName)

	if err != nil {
		return false, nil //other errors
	}

	return stat.IsDir(), nil
}
