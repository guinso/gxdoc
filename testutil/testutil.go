package testutil

import (
	"database/sql"
	"fmt"

	//explicitly include GO mysql library
	//_ "github.com/go-sql-driver/mysql"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)

var dbb *sql.DB

//GetTestDB get data base handler for test
func GetTestDB() (*sql.DB, error) {
	if dbb != nil {
		return dbb, nil
	}

	dbb, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		"root",
		"",
		"localhost",
		3306,
		"gx_doccon"))

	if err != nil {
		return nil, err
	}

	//check connection is valid or not
	if pingErr := dbb.Ping(); pingErr != nil {
		return nil, pingErr
	}

	return dbb, nil
}
