package document

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/guinso/gxschema"

	"github.com/guinso/gxdoc/util"
)

type addNewSchemaInfoItem struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type updateSchemaInfoItem struct {
	Name        string `jon:"name"`
	Description string `json:"desc"`
	IsActive    bool   `json:"isActive"`
}

var schemaRevisionPattern = regexp.MustCompile(`^document/schemas/.+/revisions/[1-9][0-9]*$`)
var schemaLatestRevPattern = regexp.MustCompile(`^document/schemas/.+$`)
var schemaDraftPattern = regexp.MustCompile(`^document/schemas/.+/draft$`)

//HandleHTTP handle HTTP request
func HandleHTTP(sanatizeURL string, w http.ResponseWriter, r *http.Request) (bool, error) {
	if sanatizeURL == "document/schema-infos" && util.IsGET(r) {
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
	} else if sanatizeURL == "document/schema-infos" && util.IsPOST(r) {
		//save new schema info
		body, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			return false, bodyErr
		}

		input := addNewSchemaInfoItem{}
		jsonErr := json.Unmarshal(body, &input)
		if jsonErr != nil {
			util.SendHTTPResponse(w, -1, "invalid input data format", "null")
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
				util.SendHTTPResponse(w, -1, err.Error(), "null")
				return true, nil
			}

			return false, err
		}
		trx.Commit()

		util.SendHTTPResponse(w, 0, "success", "null")
		return true, nil
	} else if schemaRevisionPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get specific document schema revision (return in XML format)
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]
		revision, revErr := strconv.Atoi(rawArr[4])
		if revErr != nil {
			util.SendHTTPResponse(w, -1,
				"invalid revision value (only accept integer), please check you URL",
				"null")
			return true, nil
		}

		db := util.GetDB()
		schema, schemaErr := GetSchemaByRevision(db, name, revision)
		if schemaErr != nil {
			return false, schemaErr
		}

		if schema == nil {
			util.SendHTTPResponse(w, 0, "record not found", "null")
			return true, nil
		}

		schema.ID = "" //hide ID from expose to end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			return false, xmlErr
		}

		w.Header().Set("Content-Type", "text/xml; charset=utf8")
		w.WriteHeader(200)
		w.Write([]byte(xmlStr))

		return true, nil

	} else if schemaLatestRevPattern.MatchString(sanatizeURL) && util.IsPOST(r) {
		//update document schema (input data must be XML)
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		//get input string and parse into XML
		bodyRaw, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			return false, bodyErr
		}
		bodyStr := string(bodyRaw)

		dxdoc, dxErr := gxschema.DecodeDxXML(bodyStr)
		if dxErr != nil {
			util.SendHTTPResponse(w, -1, "invalid XML: "+dxErr.Error(), "null")
			return true, nil
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			return false, trxErr
		}
		_, err := AddSchema(trx, name, dxdoc, "")
		if err != nil {
			trx.Rollback()
			return false, err
		}
		trx.Commit()

		util.SendHTTPResponse(w, 0, "update success", "null")
		return true, nil

	} else if schemaLatestRevPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get specific document schema revision (return in XML format)
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		db := util.GetDB()
		schema, schemaErr := GetSchema(db, name)
		if schemaErr != nil {
			return false, schemaErr
		}

		if schema == nil {
			util.SendHTTPResponse(w, 0, "record not found", "null")
			return true, nil
		}

		schema.ID = "" //hide ID from expose to end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			return false, xmlErr
		}

		w.Header().Set("Content-Type", "text/xml; charset=utf8")
		w.WriteHeader(200)
		w.Write([]byte(xmlStr))

		return true, nil
	} else if schemaDraftPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get draft schema
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		db := util.GetDB()
		schema, schemaErr := GetSchemaByRevision(db, name, -1)
		if schemaErr != nil {
			return false, schemaErr
		}

		if schema == nil {
			util.SendHTTPResponse(w, 0, "not record found", "null")
			return true, nil
		}

		schema.ID = "" //hide ID from end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			return false, xmlErr
		}

		w.Header().Set("Content-Type", "text/xml; charset=utf8")
		w.WriteHeader(200)
		w.Write([]byte(xmlStr))

		return true, nil
	} else if schemaDraftPattern.MatchString(sanatizeURL) && util.IsPOST(r) {
		//update draft schema
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		bodyRaw, rawErr := ioutil.ReadAll(r.Body)
		if rawErr != nil {
			return false, rawErr
		}

		gxdoc, gxErr := gxschema.DecodeDxXML(string(bodyRaw))
		if gxErr != nil {
			util.SendHTTPResponse(w, -1, "invalid input data: "+gxErr.Error(), "null")
			return true, nil
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			return false, trxErr
		}
		saveDraftErr := SaveSchemaAsDraft(trx, name, gxdoc, "")
		if saveDraftErr != nil {
			trx.Rollback()
			return false, saveDraftErr
		}
		trx.Commit()

		util.SendHTTPResponse(w, 0, "update success", "null")
		return true, nil
	} else if strings.HasPrefix(sanatizeURL, "document/schemas/") && util.IsGET(r) {
		//get single schema info
		name := sanatizeURL[17:]

		db := util.GetDB()

		schemaInfo, infoErr := GetSchemaInfo(db, name)
		if infoErr != nil {
			return false, infoErr
		}

		if schemaInfo == nil {
			util.SendHTTPResponse(w, 0, "no record", "null")
			return true, nil
		}

		util.SendHTTPResponse(w, 0, "", ExportShemaInfoToJSON(schemaInfo))
		return true, nil
	} else if strings.HasPrefix(sanatizeURL, "document/schemas/") && util.IsPOST(r) {
		//update schema info
		db := util.GetDB()

		name := sanatizeURL[17:]
		schemaInfo, infoErr := GetSchemaInfo(db, name)
		if infoErr != nil {
			return false, infoErr
		}

		if schemaInfo == nil {
			util.SendHTTPResponse(w, -1, "schema not found", "null")
			return true, nil
		}

		//parse JSON input from HTTP request
		updateItem := updateSchemaInfoItem{}
		err := util.DecodeJSON(r, &updateItem)
		if err != nil {
			util.SendHTTPResponse(w, -1,
				"unable to process user input, please check your input data format", "null")
			return true, nil
		}
		schemaInfo.Name = updateItem.Name
		schemaInfo.Description = updateItem.Description
		schemaInfo.IsActive = updateItem.IsActive

		trx, trxErr := db.Begin()
		if trxErr != nil {
			return false, trxErr
		}
		updateErr := UpdateSchemaInfo(trx, schemaInfo)
		if updateErr != nil {
			trx.Rollback()

			if _, ok := updateErr.(ErrSchemaInfoNotFound); ok {
				util.SendHTTPResponse(w, -1, "schema not exists", "null")
				return true, nil
			}

			return false, updateErr
		}
		trx.Commit()

		util.SendHTTPResponse(w, 0, "update success", "null")
		return true, nil
	} else {
		return false, nil
	}

	return true, nil
}

//ExportShemaInfoToJSON export SchemaInfo into JSON string
func ExportShemaInfoToJSON(info *SchemaInfo) string {
	return fmt.Sprintf(
		`{"name": "%s","latestRev": %d,"desc":`+
			` "%s","isActive": %t,"hasDraft": %t}`,
		info.Name, info.LatestRevision,
		info.Description, info.IsActive, info.HasDraft)
}
