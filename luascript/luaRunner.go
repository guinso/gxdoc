package luascript

import (
	"database/sql"
	"net/http"
	"time"

	lua "github.com/yuin/gopher-lua"
)

//GetLuaRunner get lua runner instance
func GetLuaRunner(db *sql.DB, r *http.Request) *lua.LState {
	//reference: https://github.com/yuin/gopher-lua

	LuaDbHandler := NewLuaMySQLHandler(db)

	L := lua.NewState()
	//x defer L.Close()

	//register GO function for Lua
	L.SetGlobal("parseDateStr", L.NewFunction(parseDateTimeString))

	L.SetGlobal("ExecSQL", L.NewFunction(LuaDbHandler.ExecSQL))
	L.SetGlobal("QuerySQL", L.NewFunction(LuaDbHandler.QuerySQL))

	return L

	// //run LUA script
	// startTime := time.Now()
	// err := L.DoString(`
	// 	-- local gg = QuerySQL("SELECT * FROM b1500f1 LIMIT 10", {}, {"num", "str", "str", "str"})
	// 	-- print(gg[3][2])

	// 	-- ExecSQL("UPDATE b1500f1 SET filename = '123abcd.csv' WHERE id = ? ", {"5"})

	// 	-- local table123 = getTableData()
	// 	-- print(table123[1][2]) -- expect get "asd"

	// 	-- local utf8String = testUtf8()
	// 	-- print(utf8String)

	// 	-- local myDt = parseDateStr("2018-05-30 14:56:04")
	// 	-- print(myDt.wday)
	// 	`)
	// if err != nil {
	// 	//panic(err)
	// 	log.Println(err.Error())
	// }
	// elapsed := time.Since(startTime)
	// log.Println(fmt.Sprintf("Lua execute in %s", elapsed))
}

//LuaParseDateTimeString parse date time string
//accepted string format is 2018-07-31 14:32:07
func parseDateTimeString(L *lua.LState) int {
	dateTimeStr := L.ToString(1) //first parameter is date time string

	t1, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
	if err != nil {
		L.RaiseError("unable to parse date time string: %s", dateTimeStr)
	}

	result := lua.LTable{}
	result.RawSetString("year", lua.LNumber(t1.Year()))
	result.RawSetString("month", lua.LNumber(t1.Month()))
	result.RawSetString("day", lua.LNumber(t1.Day()))
	result.RawSetString("yday", lua.LNumber(t1.YearDay()))
	result.RawSetString("wday", lua.LNumber(t1.Weekday()+1)) //in Lua, Sunday is 1
	result.RawSetString("hour", lua.LNumber(t1.Hour()))
	result.RawSetString("min", lua.LNumber(t1.Minute()))
	result.RawSetString("sec", lua.LNumber(t1.Second()))
	result.RawSetString("isdst", lua.LBool(false)) //is day light saving mode

	L.Push(&result)

	return 1
}
