package document

import (
	"database/sql"
	"fmt"

	"github.com/guinso/gxschema"
	"github.com/guinso/rdbmstool"
)

//GetDraftSchema get draft version of specified document schema
func GetDraftSchema(db rdbmstool.DbHandlerProxy, name string) (*gxschema.DxDoc, error) {
	return GetSchemaByRevision(db, name, -1)
}

//GetSchema get latest document schema from database
func GetSchema(db rdbmstool.DbHandlerProxy, name string) (*gxschema.DxDoc, error) {
	sqlStr := `SELECT b.xml_definition
	FROM (
	SELECT name, MAX(revision) AS latest_revision FROM doc_schema
	GROUP BY name) AS a
	JOIN doc_schema b ON a.name = b.name AND a.latest_revision = b.revision
	WHERE a.name = ?`

	row := db.QueryRow(sqlStr, name)

	var xmlDef string
	err := row.Scan(&xmlDef)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil //return NULL for record not found in database
		}

		return nil, fmt.Errorf("Failed to fetch doc_schema %s from DB: %s", name, err.Error())
	}

	//convert into DxDoc instance
	dxdoc, dxErr := gxschema.DecodeDxXML(xmlDef)
	if dxErr != nil {
		return nil, dxErr
	}

	return dxdoc, nil
}

//GetSchemaByRevision get document schema from database by revision
func GetSchemaByRevision(db rdbmstool.DbHandlerProxy, name string, revision int) (*gxschema.DxDoc, error) {
	sqlStr := `SELECT xml_definition FROM doc_schema
	WHERE name = ? AND revision = ?`

	row := db.QueryRow(sqlStr, name, revision)

	var xmlDef string
	err := row.Scan(&xmlDef)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil //return NULL for record not found
		}

		return nil, fmt.Errorf("failed to fetch doc_schema %s rev%d from database: %s",
			name, revision, err.Error())
	}

	//convert into DxDoc instance
	dxdoc, dxErr := gxschema.DecodeDxXML(xmlDef)
	if dxErr != nil {
		return nil, dxErr
	}

	return dxdoc, nil
}

//AddSchema register document schema into database
//	ErrSchemaInfoNotFound error will return if not found in doc_schema_info datatable
//RETURN:
//	int: latest revision number
func AddSchema(db rdbmstool.DbHandlerProxy, doc *gxschema.DxDoc, remark string) (int, error) {
	//check SchemaInfo is registered
	schemaInfo, infoErr := GetSchemaInfo(db, doc.Name)
	if infoErr != nil {
		return 0, infoErr
	}
	if schemaInfo == nil {
		return 0, ErrSchemaInfoNotFound{msg: doc.Name + " not found in record"}
	}

	sqlStr1 := `SELECT MAX(revision) FROM doc_schema WHERE name = ? GROUP BY name`

	//get latest revision number from database
	row := db.QueryRow(sqlStr1, doc.Name)
	var revision int
	fetchErr := row.Scan(&revision)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			revision = 0
		}

		return 0, fmt.Errorf("failed to fetch record from database: %s", fetchErr.Error())
	}

	revision++
	oriRev := doc.Revision
	doc.Revision = revision //increament by one for new revision
	xmlStr, xmlErr := doc.XML()
	doc.Revision = oriRev
	if xmlErr != nil {
		return 0, fmt.Errorf("failed to get XML definition: %s", xmlErr.Error())
	}

	_, insertErr := db.Exec(`INSERT INTO doc_schema (name,revision,xml_definition,remark) VALUES (?,?,?,?)`,
		doc.Name, revision, xmlStr, remark)

	if insertErr != nil {
		return 0, fmt.Errorf("failed to register new %s definition into database: %s",
			doc.Name, insertErr.Error())
	}

	return revision, nil
}

//SaveSchemaAsDraft save document schema as draft.
//	draft doc schema shall not affect production record
//	ErrSchemaInfoNotFound error will return if not found in doc_schema_info datatable
//NOTE: Existed draft record with same document name will be overwrite
func SaveSchemaAsDraft(db rdbmstool.DbHandlerProxy, doc *gxschema.DxDoc, remark string) error {
	xmlStr, xmlErr := doc.XML()
	if xmlErr != nil {
		return fmt.Errorf("failed convert doc schema into XML schema format: %s", xmlErr.Error())
	}

	//check SchemaInfo is registered
	schemaInfo, infoErr := GetSchemaInfo(db, doc.Name)
	if infoErr != nil {
		return infoErr
	}
	if schemaInfo == nil {
		return ErrSchemaInfoNotFound{msg: doc.Name + " not found in record"}
	}

	//check draft record is registered on database or not
	sqlStr1 := `SELECT COUNT(name) FROM doc_schema WHERE name = ? AND revision = -1 GROUP BY name`
	row := db.QueryRow(sqlStr1, doc.Name)
	var count int
	fetchErr := row.Scan(&count)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			count = 0
		} else {
			return fmt.Errorf("failed to fetch record from database: %s", fetchErr.Error())
		}
	}

	if count == 0 {
		//create a new record
		insertSQL := `INSERT INTO doc_schema (name, revision, xml_definition, remark) VALUES(?,-1,?,?)`
		_, dbErr := db.Exec(insertSQL, doc.Name, xmlStr, remark)
		if dbErr != nil {
			return fmt.Errorf("failed to save %s as draft into database: %s", doc.Name, dbErr.Error())
		}
	} else {
		//update XML_dfinition column
		updateSQL := `UPDATE doc_schema SET xml_definition = ?, remark = ? WHERE name = ? AND revision = -1`
		_, dbErr := db.Exec(updateSQL, xmlStr, remark, doc.Name)
		if dbErr != nil {
			return fmt.Errorf("failed to save %s as draft into database: %s", doc.Name, dbErr.Error())
		}
	}

	return nil
}
