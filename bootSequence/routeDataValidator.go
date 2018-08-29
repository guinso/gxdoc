package bootSequence

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/guinso/gxdoc/document"
	"github.com/guinso/gxschema"

	"github.com/guinso/gxdoc/util"
)

var dataValidatePattern = regexp.MustCompile(`^document/.+/validate$`)

//HandleDataValidationHTTP handle HTTP routing for data validation
func HandleDataValidationHTTP(sanatizeURL string, w http.ResponseWriter, r *http.Request) bool {
	if !dataValidatePattern.MatchString(sanatizeURL) || !util.IsPOST(r) {
		return false //URL pattern not match
	}

	docSchemaName := strings.Split(sanatizeURL, "/")[1]
	docSchema, schemaErr := document.GetSchema(util.GetDB(), docSchemaName)
	if schemaErr != nil {
		//TODO: log error
		util.LogError(schemaErr)
		util.SendHTTPServerErrorJSON(w)
	}

	inputStr, inputErr := util.GetHTTPRequestBody(r)
	if inputErr != nil {
		//TODO: log error
		util.LogError(inputErr)
		util.SendHTTPServerErrorJSON(w)
	}

	//validate data in JSON or XML format
	dataTypeRaw := strings.Split(r.Header.Get("Content-Type"), ";")[0]
	if strings.Compare("application/json", dataTypeRaw) == 0 {
		invalid := gxschema.ValidateDataFromJSON(inputStr, docSchema)
		if invalid == nil {
			util.SendHTTPResponseJSON(w, `{"isValid":true, "message":""}`)
		} else {
			util.SendHTTPResponseJSON(w,
				fmt.Sprintf(`{"isValid":false, "message":"%s"}`, invalid.Error()))
		}
	} else if strings.Compare("text/xml", dataTypeRaw) == 0 {
		invalid := gxschema.ValidateDataFromXML(inputStr, docSchema)
		if invalid == nil {
			util.SendHTTPResponseJSON(w, `{"isValid":true, "message":""}`)
		} else {
			util.SendHTTPResponseJSON(w,
				fmt.Sprintf(`{"isValid":false, "message":"%s"}`, invalid.Error()))
		}
	} else {
		util.SendHTTPClientErrorJSON(w, 400, -1, "input data type only accept either JSON nor XML")
	}

	return true
}
