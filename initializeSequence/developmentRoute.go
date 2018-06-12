package initializeSequence

import (
	"log"
	"net/http"

	"github.com/guinso/gxdoc/routing"

	"github.com/guinso/gxdoc/document"
)

//routeDevelopmentAPI handle HTTP REST API (development)
func routeDevelopmentAPI(w http.ResponseWriter, r *http.Request, url string) {
	//TODO: how to check requstor ID and determine her authority?

	done, docErr := document.HandleHTTP(url, w, r)
	if done {
		return
	} else if docErr != nil {
		if _, ok := docErr.(routing.ErrInvalidInputData); ok {
			http.Error(w, "invalid input data", 400)
		} else if _, ok := docErr.(routing.ErrNotAuthorize); ok {
			http.Error(w, "unauthorize", 401)
		} else if _, ok := docErr.(routing.ErrNotAuthenticate); ok {
			http.Error(w, "unauthorize", 401)
		} else {
			log.Println(docErr.Error())
			http.Error(w, "internal error", 500)
		}
	}

	http.Error(w, "path not found", 404)
}
