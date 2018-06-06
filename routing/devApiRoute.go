package routing

import "net/http"

func routeDevPath(w http.ResponseWriter, r *http.Request, trimURL string) {
	//TODO: handle development API
	handleErrorCode(404, "Path not found.", w)
}
