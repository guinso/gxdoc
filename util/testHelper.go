package util

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//HTTPMsg return message format for goweb REST API
type HTTPMsg struct {
	StatusCode    int         `json:"statusCode"`
	StatusMessage string      `json:"statusMsg"`
	Response      interface{} `json:"response,omitempty"`
}

const (
	TestDatabaseName = "gx_doccon"
)

var dbb *sql.DB

//GetTestDB get database handler for unit test
//WARNING: don't use it in source code other than unit test!
//       : make sure it is pointed to non-production database for testing purpose
func GetTestDB() *sql.DB {

	if dbb == nil {
		dbTest, dbErr := sql.Open("mysql", fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8",
			"root",            //unit test username
			"",                //unit test password
			"localhost",       //unit test server location
			3306,              //unit test database port number
			TestDatabaseName)) //unit test database name

		if dbErr != nil {
			panic(dbErr)
		}

		dbb = dbTest
	}

	return dbb
}

//RestRequestForTest tool to send REST request to goweb server
func RestRequestForTest(
	request *http.Request,
	responseObj interface{}) (*HTTPMsg, *http.Response, error) {

	client := http.Client{}
	response, resErr := client.Do(request)
	if resErr != nil {
		return nil, nil, resErr
	}
	//defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, nil, fmt.Errorf("Internal server error: %d", response.StatusCode)
	}

	var rawMsg json.RawMessage
	returnObj := HTTPMsg{
		Response: &rawMsg,
	}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&returnObj); err != nil {
		return nil, nil, err
	}

	//decode custom field
	if customDecodeErr := json.Unmarshal(rawMsg, responseObj); customDecodeErr != nil {
		return nil, nil, errors.New("Failed to decode response body: " + customDecodeErr.Error())
	}

	return &returnObj, response, nil
}
