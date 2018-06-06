package luascript

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/yuin/gopher-lua"
)

//LuaMySQLHandler Lua MySQl handler
type LuaMySQLHandler struct {
	//DbHandler MySQL database handler
	DbHandler *sql.DB
}

//NewLuaMySQLHandler create a new instance of Lua MySQL handler
func NewLuaMySQLHandler(dbHandler *sql.DB) *LuaMySQLHandler {
	return &LuaMySQLHandler{DbHandler: dbHandler}
}

//RegisterLuaBindings register nessesary SQL functions and datatype to Lua virtual enviroment
func (sqlHandler *LuaMySQLHandler) RegisterLuaBindings(L *lua.LState) {
	L.SetGlobal("SqlExec", L.NewFunction(sqlHandler.ExecSQL))
	L.SetGlobal("SqlQuery", L.NewFunction(sqlHandler.QuerySQL))
}

//QuerySQL query SQL statement
func (sqlHandler *LuaMySQLHandler) QuerySQL(L *lua.LState) int {
	sqlStr := L.ToString(1)      //get first argument as SQL statement
	inputParams := L.ToTable(2)  //get 2nd argument as SQL parameters
	outputFormat := L.ToTable(3) //get 3rd argument as output row's data format

	stmt, err := sqlHandler.DbHandler.Prepare(sqlStr)
	if err != nil {
		log.Println(err.Error())
		L.RaiseError(err.Error())
		return 0
	}

	//convert parameters from lua
	//paramLen := inputParams.Len()
	params := make([]interface{}, 0)
	var tmpType lua.LValueType
	hasError := false
	if inputParams != nil && inputParams.Len() > 0 {
		inputParams.ForEach(func(key lua.LValue, value lua.LValue) {
			tmpType = value.Type()

			if tmpType == lua.LTNumber {
				tmpFloat, floatErr := strconv.ParseFloat(value.String(), 64)
				if floatErr != nil {
					log.Println(floatErr.Error())
					L.RaiseError(floatErr.Error())
					hasError = true
					return
				}
				params = append(params, tmpFloat)

			} else if tmpType == lua.LTBool {
				tmpBool, boolfloatErr := strconv.ParseBool(value.String())
				if boolfloatErr != nil {
					log.Println(boolfloatErr.Error())
					L.RaiseError(boolfloatErr.Error())
					hasError = true
					return
				}
				params = append(params, tmpBool)

			} else if tmpType == lua.LTString {
				params = append(params, value.String())
			} else {
				errMsg := fmt.Sprintf("Parameter type of %s is not supported as SQL parameter", value.Type().String())
				log.Println(errMsg)
				L.RaiseError(errMsg)
				hasError = true
				return
			}
		})
	}

	if hasError {
		return 0
	}

	rows, queryErr := stmt.Query(params...)
	if queryErr != nil {
		log.Println("QuerySQL: failed to query SQL: " + queryErr.Error())
		L.RaiseError(queryErr.Error())
	}

	result := lua.LTable{}
	defer rows.Close()
	for rows.Next() {
		tmpRow, tmpErr := sqlHandler.prepareInterfaceArray(outputFormat)
		if tmpErr != nil {
			log.Println("QuerySQL: failed to prepare output row: " + tmpErr.Error())
			L.RaiseError(tmpErr.Error())
			return 0
		}

		//fetch value from DB result row
		scanErr := rows.Scan(tmpRow...)
		if scanErr != nil {
			log.Println("QuerySQL: failed to fetch row: " + scanErr.Error())
			L.RaiseError(scanErr.Error())
			return 0
		}

		//convert row value into LTable row
		row := lua.LTable{}
		for i := 0; i < len(tmpRow); i++ {
			if x, ok := tmpRow[i].(*sql.NullBool); ok == true {
				if !x.Valid {
					//log.Println("try push nil (bool)")
					row.Append(lua.LNil)
				} else {
					//log.Println("try push bool")
					row.Append(lua.LBool(x.Bool))
				}
			} else if x, ok := tmpRow[i].(*sql.NullString); ok == true {
				if !x.Valid {
					//log.Println("try push nil (string)")
					row.Append(lua.LNil)
				} else {
					//log.Println("try push string")
					row.Append(lua.LString(x.String))
				}
			} else if x, ok := tmpRow[i].(*sql.NullFloat64); ok == true {
				if !x.Valid {
					//log.Println("try push nil (float64)")
					row.Append(lua.LNil)
				} else {
					//log.Println("try push float64")
					row.Append(lua.LNumber(x.Float64))
				}
			} else {
				errMsg := fmt.Sprintf(
					"QuerySQL: Result column type %s at index %d is not supported to convert into lua table value",
					reflect.TypeOf(tmpRow[i]), i)
				log.Println(errMsg)
				L.RaiseError(errMsg)
				return 0
			}
		}

		result.Append(&row)
	}

	L.Push(&result)

	return 1 //indicate there is one result is return to user
}

func (sqlHandler *LuaMySQLHandler) prepareInterfaceArray(dataFormat *lua.LTable) ([]interface{}, error) {
	result := make([]interface{}, 0)

	var conversionErr error
	var tmpType string
	dataFormat.ForEach(func(key lua.LValue, value lua.LValue) {
		if conversionErr != nil {
			return
		}

		tmpType = strings.ToLower(value.String())

		if strings.Compare(tmpType, "num") == 0 {
			result = append(result, &sql.NullFloat64{})
		} else if strings.Compare(tmpType, "str") == 0 {
			result = append(result, &sql.NullString{})
		} else if strings.Compare(tmpType, "bool") == 0 {
			result = append(result, &sql.NullBool{})
		} else {
			conversionErr = fmt.Errorf("Type %s is not supported to prepare SQL parameters", tmpType)
		}
	})

	if conversionErr != nil {
		return nil, conversionErr
	}

	return result, nil
}

//ExecSQL execute SQL statement
//Lua example: SqlExec("INSERT INTO account VALUES ?, ?", {1, "John"})
//    Param 1: sql statement
//    Param 2: sql input parameter(s)
func (sqlHandler *LuaMySQLHandler) ExecSQL(L *lua.LState) int {
	sqlStr := L.ToString(1)     //get first argument as string (SQL statement)
	inputParams := L.ToTable(2) //get 2nd argument as input parameters

	stmt, err := sqlHandler.DbHandler.Prepare(sqlStr)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	//convert parameters from lua
	//paramLen := inputParams.Len()
	params := make([]interface{}, 0)
	var tmpType lua.LValueType
	if inputParams != nil && inputParams.Len() > 0 {
		inputParams.ForEach(func(key lua.LValue, value lua.LValue) {
			tmpType = value.Type()

			if tmpType == lua.LTNumber {
				tmpFloat, floatErr := strconv.ParseFloat(value.String(), 64)
				if floatErr != nil {
					L.RaiseError(floatErr.Error())
					return
				}
				params = append(params, tmpFloat)

			} else if tmpType == lua.LTBool {
				tmpBool, boolfloatErr := strconv.ParseBool(value.String())
				if boolfloatErr != nil {
					L.RaiseError(boolfloatErr.Error())
					return
				}
				params = append(params, tmpBool)

			} else if tmpType == lua.LTString {
				params = append(params, value.String())
			} else {
				L.RaiseError("Parameter type of %s is not supported as SQL parameter", value.Type().String())
			}
		})
	}

	_, execErr := stmt.Exec(params...)
	if execErr != nil {
		L.RaiseError("ExecSQL failed to execute SQL: %s, \r\nerror message: %s", sqlStr, execErr.Error())
	}

	return 0 //indicate there is no result is return to user
}
