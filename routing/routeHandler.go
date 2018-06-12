package routing

import "net/http"

//RouteHandler handle HTTP request interface
type RouteHandler interface {
	HandleHTTP(sanatizeURL string, w http.ResponseWriter, r *http.Request) error
}
