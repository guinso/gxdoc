package luascript

import (
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

//LuaHTTPHandler Lua HTTP handler
type LuaHTTPHandler struct {
	Request *http.Request
	Writer  *http.ResponseWriter
}

//RegisterLuaBindings register golang HTTP functionality to Lua runner
func (handler *LuaHTTPHandler) RegisterLuaBindings(L *lua.LState) {

}

//GetCookieValue get value from cookie by providing key
func (handler *LuaHTTPHandler) GetCookieValue(L *lua.LState) int {
	cookieKey := L.ToString(1) //get 1st parameter as cookie key

	cookie, cookieErr := handler.Request.Cookie(cookieKey)
	if cookieErr != nil {
		//L.RaiseError(cookieErr.Error())
		L.Push(lua.LNil) //cookie not found
	}

	if cookie == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(cookie.Value))
	}

	return 1
}

//SetCookieValue set string value into cookie
func (handler *LuaHTTPHandler) SetCookieValue(L *lua.LState) int {
	cookieKey := L.ToString(1) //get 1st parameter as cookie key
	value := L.ToString(2)     //get 2nd parameter as new cookie value

	cookie, cookieErr := handler.Request.Cookie(cookieKey)
	if cookieErr != nil {
		L.RaiseError(cookieErr.Error())
	}

	cookie.Value = value

	return 0
}
