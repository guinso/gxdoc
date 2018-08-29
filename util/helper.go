package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/guinso/rdbmstool"

	//explicitly include GO mysql library
	//_ "github.com/go-sql-driver/mysql"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)

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

//SendHTTPResponseXML send HTTP response in XML format
func SendHTTPResponseXML(w http.ResponseWriter, xml string) {
	w.Header().Set("Content-Type", "text/xml; charset=utf8")
	w.WriteHeader(200)
	w.Write([]byte(xml))
}

//SendHTTPResponseJSON send HTTP response in JSON format
func SendHTTPResponseJSON(w http.ResponseWriter, json string) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(200)
	w.Write([]byte(json))
}

//SendHTTPClientErrorJSON send HTTP error response to client
//due to client request is rejected at server side (in JSON format)
//
//	httpCode is HTTP status code, please use value in range 400-499, use 400 if has no idea
//	errorCode is application's itself's error code, defined by developer, use -1 if has no idea
//	errorMessage is description to describe errors
func SendHTTPClientErrorJSON(w http.ResponseWriter, httpCode int, errorCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(httpCode)
	w.Write([]byte(fmt.Sprintf(`{"errorCode":%d, "errorMessage":"%s"}`, errorCode, errorMessage)))
}

//SendHTTPServerErrorJSON send HTTP error response to client
//due to server side encounter unexpected error (in JSON format)
func SendHTTPServerErrorJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(500)
	w.Write([]byte("{msg:\"Encounter internal server error\"}"))
}

//IsPOST check request HTTP is POST method
func IsPOST(r *http.Request) bool {
	return strings.Compare(r.Method, "POST") == 0
}

//IsGET check request HTTP is GET method
func IsGET(r *http.Request) bool {
	return strings.Compare(r.Method, "GET") == 0
}

//IsPUT check request HTTP is PUT method
func IsPUT(r *http.Request) bool {
	return strings.Compare(r.Method, "PUT") == 0
}

//IsDELETE check request HTTP is DELETE method
func IsDELETE(r *http.Request) bool {
	return strings.Compare(r.Method, "DELETE") == 0
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

//GetDBColumnMaxInt get maximum integer value from particular data table's column
func GetDBColumnMaxInt(db rdbmstool.DbHandlerProxy, tableName string, columnName string) (int, error) {
	var tmpInt sql.NullInt64
	sqlStr := fmt.Sprintf("SELECT MAX(%s) FROM %s", columnName, tableName)
	row := db.QueryRow(sqlStr)
	scanErr := row.Scan(&tmpInt)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return 0, nil
		}

		return 0, scanErr
	}

	if tmpInt.Valid {
		return int(tmpInt.Int64), nil
	}

	return 0, nil
}
