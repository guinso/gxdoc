package routing

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/guinso/gxdoc/document"
	"github.com/guinso/gxdoc/util"
)

//routeDevelopmentAPI handle HTTP REST API (development)
func routeDevelopmentAPI(w http.ResponseWriter, r *http.Request, url string) {
	//util.SendHTTPResponse(w, 0, "ok", "{}")

	//TODO: how to check requstor ID and determine her authority?

	if strings.Compare(url, "document/schema-info") == 0 {
		if util.IsGET(r) {
			//TODO: get all document schema info
			results, err := document.GetAllSchemaInfo(util.GetDB())
			if err != nil {
				handleErrorCode(500, "internal server error", w)
				return
			}

			byteArr, byteErr := json.Marshal(results)
			if byteErr != nil {
				handleErrorCode(500, "internal server error", w)
				return
			}

			util.SendHTTPResponse(w, 0, "ok", string(byteArr))
			return
		} else if util.IsPOST(r) {
			//TODO: register a new document schema info
		}
	}

	// if strings.HasPrefix(url, "document/schema-info/") {
	// 	schemaID := url[21:]
	// 	if util.IsGET(r) {
	// 		//TODO: get specified document schema info
	// 	} else if util.IsPOST(r) {
	// 		//TODO: update specified document schema info
	// 	}
	// }

	// if strings.HasPrefix(url, "document/schema/") {
	// 	//handle document schema
	// 	schemaID := url[16:]
	// }

	// if strings.HasPrefix(url, "document/record/") {
	// 	//handle document data
	// 	schemaID := url[16:]
	// }

	handleErrorCode(404, "Path not found.", w)
}
