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
func HandleDocSchemaHTTP(sanatizeURL string, w http.ResponseWriter, r *http.Request) (bool, error) {
	if sanatizeURL == "document/schema-infos" && util.IsGET(r) {
		//get all document schema info
		results, err := document.GetAllSchemaInfo(util.GetDB())
		if err != nil {
			return false, err
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

		return true, nil
	} else if sanatizeURL == "document/schema-infos" && util.IsPOST(r) {
		//save new schema info
		body, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			return false, bodyErr
		}

		input := addNewSchemaInfoItem{}
		jsonErr := json.Unmarshal(body, &input)
		if jsonErr != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1, "invalid input data format")
			return true, nil
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			return false, trxErr
		}
		err := document.AddSchemaInfo(db, input.Name, input.Description)
		if err != nil {
			trx.Rollback()
			if _, ok := err.(document.ErrSchemaInfoAlreadyExists); ok {
				util.SendHTTPClientErrorJSON(w, 400, -1, err.Error())
				return true, nil
			}

			return false, err
		}
		trx.Commit()

		util.SendHTTPResponseJSON(w, "{}")

		return true, nil
	} else if schemaRevisionPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get specific document schema revision (return in XML format)
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]
		revision, revErr := strconv.Atoi(rawArr[4])
		if revErr != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1,
				"invalid revision value (only accept integer), please check you URL")
			return true, nil
		}

		db := util.GetDB()
		schema, schemaErr := document.GetSchemaByRevision(db, name, revision)
		if schemaErr != nil {
			return false, schemaErr
		}

		if schema == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "record not found")
			return true, nil
		}

		schema.ID = "" //hide ID from expose to end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			return false, xmlErr
		}

		util.SendHTTPResponseXML(w, xmlStr)

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

		dxdoc, dxErr := gxschema.ParseSchemaFromXML(bodyStr)
		if dxErr != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1, "invalid XML: "+dxErr.Error())
			return true, nil
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			return false, trxErr
		}
		_, err := document.AddSchema(trx, name, dxdoc, "")
		if err != nil {
			trx.Rollback()
			return false, err
		}
		trx.Commit()

		util.SendHTTPResponseJSON(w, "{}")
		return true, nil

	} else if schemaLatestRevPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get specific document schema revision (return in XML format)
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		db := util.GetDB()
		schema, schemaErr := document.GetSchema(db, name)
		if schemaErr != nil {
			return false, schemaErr
		}

		if schema == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "record not found")
			return true, nil
		}

		schema.ID = "" //hide ID from expose to end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			return false, xmlErr
		}

		//TODO: include XSD reference as well
		util.SendHTTPResponseXML(w, xmlStr)

		return true, nil
	} else if schemaDraftPattern.MatchString(sanatizeURL) && util.IsGET(r) {
		//get draft schema
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		db := util.GetDB()
		schema, schemaErr := document.GetSchemaByRevision(db, name, -1)
		if schemaErr != nil {
			return false, schemaErr
		}

		if schema == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "not record found")
			return true, nil
		}

		schema.ID = "" //hide ID from end user
		xmlStr, xmlErr := schema.XML()
		if xmlErr != nil {
			return false, xmlErr
		}

		util.SendHTTPResponseXML(w, xmlStr)

		return true, nil
	} else if schemaDraftPattern.MatchString(sanatizeURL) && util.IsPOST(r) {
		//update draft schema
		rawArr := strings.Split(sanatizeURL, "/")
		name := rawArr[2]

		bodyRaw, rawErr := ioutil.ReadAll(r.Body)
		if rawErr != nil {
			return false, rawErr
		}

		gxdoc, gxErr := gxschema.ParseSchemaFromXML(string(bodyRaw))
		if gxErr != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1, "invalid input data: "+gxErr.Error())
			return true, nil
		}

		db := util.GetDB()
		trx, trxErr := db.Begin()
		if trxErr != nil {
			return false, trxErr
		}
		saveDraftErr := document.SaveSchemaAsDraft(trx, name, gxdoc, "")
		if saveDraftErr != nil {
			trx.Rollback()
			return false, saveDraftErr
		}
		trx.Commit()

		util.SendHTTPResponseJSON(w, "{}")
		return true, nil
	} else if strings.HasPrefix(sanatizeURL, "document/schemas/") && util.IsGET(r) {
		//get single schema info
		name := sanatizeURL[17:]

		db := util.GetDB()

		schemaInfo, infoErr := document.GetSchemaInfo(db, name)
		if infoErr != nil {
			return false, infoErr
		}

		if schemaInfo == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "no record")
			return true, nil
		}

		util.SendHTTPResponseJSON(w, schemaInfo.JSON())
		return true, nil
	} else if strings.HasPrefix(sanatizeURL, "document/schemas/") && util.IsPOST(r) {
		//update schema info
		db := util.GetDB()

		name := sanatizeURL[17:]
		schemaInfo, infoErr := document.GetSchemaInfo(db, name)
		if infoErr != nil {
			return false, infoErr
		}

		if schemaInfo == nil {
			util.SendHTTPClientErrorJSON(w, 404, -1, "schema not found")
			return true, nil
		}

		//parse JSON input from HTTP request
		updateItem := updateSchemaInfoItem{}
		err := util.DecodeJSON(r, &updateItem)
		if err != nil {
			util.SendHTTPClientErrorJSON(w, 400, -1,
				"unable to process user input, please check your input data format")
			return true, nil
		}
		schemaInfo.Name = updateItem.Name
		schemaInfo.Description = updateItem.Description
		schemaInfo.IsActive = updateItem.IsActive

		trx, trxErr := db.Begin()
		if trxErr != nil {
			return false, trxErr
		}
		updateErr := document.UpdateSchemaInfo(trx, schemaInfo)
		if updateErr != nil {
			trx.Rollback()

			if _, ok := updateErr.(document.ErrSchemaInfoNotFound); ok {
				util.SendHTTPClientErrorJSON(w, 404, -1, "schema not exists")
				return true, nil
			}

			return false, updateErr
		}
		trx.Commit()

		util.SendHTTPResponseJSON(w, "{}")
		return true, nil
	}

	return false, nil
}
