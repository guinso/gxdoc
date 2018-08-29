package bootSequence

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/guinso/gxdoc/document"
	"github.com/guinso/gxdoc/util"
	"github.com/guinso/gxschema"
)

//addNewSchemaInfoItem add schema info data type
type addNewSchemaInfoItem struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}

//updateSchemaInfoItem update schema info data type
type updateSchemaInfoItem struct {
	Name        string `jon:"name"`
	Description string `json:"desc"`
	IsActive    bool   `json:"isActive"`
}

var schemaRevisionPattern = regexp.MustCompile(`^document/schemas/.+/revisions/[1-9][0-9]*$`)
var schemaLatestRevPattern = regexp.MustCompile(`^document/schemas/.+$`)
var schemaDraftPattern = regexp.MustCompile(`^document/schemas/.+/draft$`)

//HandleDocSchemaHTTP handle HTTP request
func HandleDocSchemaHTTP(sanatizeURL string, w http.ResponseWriter, r *http.Request) bool {
	if sanatizeURL == "document/schema-infos" && util.IsGET(r) {
		//get all document schema info
		results, err := document.GetAllSchemaInfo(util.GetDB())
		if err != nil {
			//TODO: log error message
			util.LogError(err)
			util.SendHTTPServerErrorJSON(w)
			return true
		}

		var result string
		for index, item := range results {
			if index == 0 {
				result = item.JSON()
			} else {
				result = result + "," + item.JSON()
			}
		}
		result = "[" + result + "]"

		util.SendHTTPResponseJSON(w, result)

		return true
	} else if sanatizeURL == "document/schema-infos" && util.IsPOST(r) {
		//save new schema info
		body, bodyErr := util.GetHTTPRequestBody(r)
		if bodyErr != nil {
			util.LogError(bodyErr)
			util.SendHTTPServerErrorJSON(w)
			return true
		}

		input := addNewSchemaInfoItem{}
		jsonErr := json.Unmarshal([]byte(body), &input)
		if jsonErr != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1, "invalid input data format")
			return true
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			util.LogError(trxErr)
			util.SendHTTPServerErrorJSON(w)
			return true
		}
		err := document.AddSchemaInfo(db, input.Name, input.Description)
		if err != nil {
			trx.Rollback()
			if _, ok := err.(document.ErrSchemaInfoAlreadyExists); ok {
				util.SendHTTPClientErrorJSON(w, 400, -1, err.Error())
				return true
			}

			util.LogError(err)
			util.SendHTTPServerErrorJSON(w)

			return true
		}
		trx.Commit()

		util.SendHTTPResponseJSON(w, "{}")

		return true
	} else if schemaRevisionPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get specific document schema revision (return in XML format)
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]
		revision, revErr := strconv.Atoi(rawArr[4])
		if revErr != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1,
				"invalid revision value (only accept integer), please check you URL")
			return true
		}

		db := util.GetDB()
		schema, schemaErr := document.GetSchemaByRevision(db, name, revision)
		if schemaErr != nil {
			util.LogError(schemaErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}

		if schema == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "record not found")
			return true
		}

		schema.ID = "" //hide ID from expose to end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			util.LogError(xmlErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}

		util.SendHTTPResponseXML(w, xmlStr)

		return true

	} else if schemaLatestRevPattern.MatchString(sanatizeURL) && util.IsPOST(r) {
		//update document schema (input data must be XML)
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		//get input string and parse into XML
		bodyStr, bodyErr := util.GetHTTPRequestBody(r)
		if bodyErr != nil {
			util.LogError(bodyErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}

		dxdoc, dxErr := gxschema.ParseSchemaFromXML(bodyStr)
		if dxErr != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1, "invalid XML: "+dxErr.Error())
			return true
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			util.LogError(trxErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}
		_, err := document.AddSchema(trx, name, dxdoc, "")
		if err != nil {
			trx.Rollback()

			util.LogError(err)
			util.SendHTTPServerErrorJSON(w)
			return true
		}
		trx.Commit()

		util.SendHTTPResponseJSON(w, "{}")
		return true

	} else if schemaLatestRevPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get specific document schema revision (return in XML format)
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		db := util.GetDB()
		schema, schemaErr := document.GetSchema(db, name)
		if schemaErr != nil {
			util.LogError(schemaErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}

		if schema == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "record not found")
			return true
		}

		schema.ID = "" //hide ID from expose to end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			util.LogError(xmlErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}

		//TODO: include XSD reference as well
		util.SendHTTPResponseXML(w, xmlStr)

		return true
	} else if schemaDraftPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get draft schema
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		db := util.GetDB()
		schema, schemaErr := document.GetSchemaByRevision(db, name, -1)
		if schemaErr != nil {
			util.LogError(schemaErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}

		if schema == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "not record found")
			return true
		}

		schema.ID = "" //hide ID from end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			util.LogError(xmlErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}

		util.SendHTTPResponseXML(w, xmlStr)

		return true
	} else if schemaDraftPattern.MatchString(sanatizeURL) && util.IsPOST(r) {
		//update draft schema
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		bodyRaw, rawErr := ioutil.ReadAll(r.Body)
		if rawErr != nil {
			util.LogError(rawErr)
			util.SendHTTPServerErrorJSON(w)

			return true
		}

		gxdoc, gxErr := gxschema.ParseSchemaFromXML(string(bodyRaw))
		if gxErr != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1, "invalid input data: "+gxErr.Error())
			return true
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			util.LogError(trxErr)
			util.SendHTTPServerErrorJSON(w)
			return true
		}
		saveDraftErr := document.SaveSchemaAsDraft(trx, name, gxdoc, "")
		if saveDraftErr != nil {
			trx.Rollback()

			util.LogError(saveDraftErr)
			util.SendHTTPServerErrorJSON(w)
			return true
		}
		trx.Commit()

		util.SendHTTPResponseJSON(w, "{}")
		return true
	} else if strings.HasPrefix(sanatizeURL, "document/schemas/") && util.IsGET(r) {
		//get single schema info
		name := sanatizeURL[17:]

		db := util.GetDB()

		schemaInfo, infoErr := document.GetSchemaInfo(db, name)
		if infoErr != nil {
			util.LogError(infoErr)
			util.SendHTTPServerErrorJSON(w)
			return true
		}

		if schemaInfo == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "no record")
			return true
		}

		util.SendHTTPResponseJSON(w, schemaInfo.JSON())
		return true
	} else if strings.HasPrefix(sanatizeURL, "document/schemas/") && util.IsPOST(r) {
		//update schema info
		db := util.GetDB()

		name := sanatizeURL[17:]
		schemaInfo, infoErr := document.GetSchemaInfo(db, name)
		if infoErr != nil {
			util.LogError(infoErr)
			util.SendHTTPServerErrorJSON(w)
			return true
		}

		if schemaInfo == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "schema not found")
			return true
		}

		//parse JSON input from HTTP request
		updateItem := updateSchemaInfoItem{}
		err := util.DecodeJSON(r, &updateItem)
		if err != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1,
				"unable to process user input, please check your input data format")
			return true
		}
		schemaInfo.Name = updateItem.Name
		schemaInfo.Description = updateItem.Description
		schemaInfo.IsActive = updateItem.IsActive

		trx, trxErr := db.Begin()
		if trxErr != nil {
			util.LogError(trxErr)
			util.SendHTTPServerErrorJSON(w)
			return true
		}
		updateErr := document.UpdateSchemaInfo(trx, schemaInfo)
		if updateErr != nil {
			trx.Rollback()

			if _, ok := updateErr.(document.ErrSchemaInfoNotFound); ok {
				util.SendHTTPClientErrorJSON(w, 404, -1, "schema not exists")
				return true
			}

			util.LogError(updateErr)
			util.SendHTTPServerErrorJSON(w)
			return true
		}
		trx.Commit()

		util.SendHTTPResponseJSON(w, "{}")
		return true
	}

	return false
}
