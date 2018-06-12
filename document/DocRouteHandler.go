package document

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/guinso/gxdoc/util"
)

type addNewSchemaInfo struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}

//HandleHTTP handle HTTP request
func HandleHTTP(sanatizeURL string, w http.ResponseWriter, r *http.Request) (bool, error) {
	if sanatizeURL == "document/schema-info" && util.IsGET(r) {
		//get all document schema info
		results, err := GetAllSchemaInfo(util.GetDB())
		if err != nil {
			return false, err
		}

		var result string
		for index, item := range results {
			if index == 0 {
				result = ExportShemaInfoToJSON(&item)
			} else {
				result = result + "," + ExportShemaInfoToJSON(&item)
			}
		}
		result = "[" + result + "]"

		util.SendHTTPResponse(w, 0, "ok", result)
	} else if sanatizeURL == "document/schema-info" && util.IsPOST(r) {
		//TODO: save new schema info
		body, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			return false, bodyErr
		}

		input := addNewSchemaInfo{}
		jsonErr := json.Unmarshal(body, &input)
		if jsonErr != nil {
			util.SendHTTPResponse(w, -1, "invalid input data format", "{}")
			return true, nil
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			return false, trxErr
		}
		err := AddSchemaInfo(db, input.Name, input.Description)
		if err != nil {
			trx.Rollback()
			if _, ok := err.(ErrSchemaInfoAlreadyExists); ok {
				util.SendHTTPResponse(w, -1, err.Error(), "{}")
				return true, nil
			}

			return false, err
		}
		trx.Commit()

		util.SendHTTPResponse(w, 0, "success", "{}")
		return true, nil
	} else if strings.HasPrefix(sanatizeURL, "document/schema-info/") && util.IsPUT(r) {
		//TODO: update schema info
	} else if strings.HasPrefix(sanatizeURL, "document/schema/") && util.IsGET(r) {
		//TODO: get document schema (return in XML format)
	} else if strings.HasPrefix(sanatizeURL, "document/schema/") && util.IsPUT(r) {
		//TODO: update document schema (input data must be XML)
	} else {
		return false, nil
	}

	return true, nil
}

//ExportShemaInfoToJSON export SchemaInfo into JSON string
func ExportShemaInfoToJSON(info *SchemaInfo) string {
	return fmt.Sprintf(
		`{"id": "%s","name": "%s","latestRev": %d,"desc":`+
			` "%s","isActive": %t,"hasDraft": %t}`,
		info.ID, info.Name, info.LatestRevision,
		info.Description, info.IsActive, info.HasDraft)
}
