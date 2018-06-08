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
	sqlStr := `SELECT b.xml_definition, a.latest_revision, a.schema_id 
	FROM (
	SELECT schema_id, MAX(revision) AS latest_revision FROM doc_schema_revision
	GROUP BY schema_id) AS a
	JOIN doc_schema_revision b ON a.schema_id = b.schema_id AND a.latest_revision = b.revision
	JOIN doc_schema c ON a.schema_id = c.id
	WHERE c.name = ?`

	row := db.QueryRow(sqlStr, name)

	var xmlDef string
	var tmpRev int
	var tmpID int
	err := row.Scan(&xmlDef, &tmpRev, &tmpID)
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

	//make sure doc schema tally with doc_schema_revision
	dxdoc.Name = name
	dxdoc.Revision = tmpRev
	dxdoc.ID = tmpID

	return dxdoc, nil
}

//GetSchemaByRevision get document schema from database by revision
func GetSchemaByRevision(db rdbmstool.DbHandlerProxy, name string, revision int) (*gxschema.DxDoc, error) {
	sqlStr := `SELECT a.xml_definition, a.schema_id FROM doc_schema_revision a
	JOIN doc_schema b ON a.schema_id = b.id
	WHERE b.name = ? AND a.revision = ?`

	row := db.QueryRow(sqlStr, name, revision)

	var xmlDef string
	var tmpID int
	err := row.Scan(&xmlDef, &tmpID)

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

	//make sure doc schema tally with doc_schema_revision
	dxdoc.Name = name
	dxdoc.Revision = revision
	dxdoc.ID = tmpID

	return dxdoc, nil
}

//GetSchemaByID get document schema from database by schema_id
func GetSchemaByID(db rdbmstool.DbHandlerProxy, id int) (*gxschema.DxDoc, error) {
	sqlStr := `SELECT b.xml_definition, a.latest_revision, c.name
	FROM (
	SELECT schema_id, MAX(revision) AS latest_revision FROM doc_schema_revision
	GROUP BY schema_id) AS a
	JOIN doc_schema_revision b ON a.schema_id = b.schema_id AND a.latest_revision = b.revision
	JOIN doc_schema c ON a.schema_id = c.id
	WHERE a.schema_id = ?`

	row := db.QueryRow(sqlStr, id)

	var xmlDef string
	var tmpRev int
	var tmpName string
	err := row.Scan(&xmlDef, &tmpRev, &tmpName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil //return NULL for record not found in database
		}

		return nil, fmt.Errorf("Failed to fetch doc_schema ID %d from DB : %s", id, err.Error())
	}

	//convert into DxDoc instance
	dxdoc, dxErr := gxschema.DecodeDxXML(xmlDef)
	if dxErr != nil {
		return nil, dxErr
	}

	//make sure doc schema tally with doc_schema_revision
	dxdoc.Name = tmpName
	dxdoc.Revision = tmpRev
	dxdoc.ID = id

	return dxdoc, nil
}

//AddSchema register document schema into database
//RETURN:
//	int: latest revision number
//NOTE: ErrSchemaInfoNotFound error will return if document not register yet in doc_schema_info datatable
func AddSchema(db rdbmstool.DbHandlerProxy, schemaName string, doc *gxschema.DxDoc, remark string) (int, error) {
	//check SchemaInfo is registered
	schemaInfo, infoErr := GetSchemaInfo(db, schemaName)
	if infoErr != nil {
		return 0, infoErr
	}
	if schemaInfo == nil {
		return 0, ErrSchemaInfoNotFound{msg: schemaName + " doc schema not found in record"}
	}

	sqlStr1 := `SELECT MAX(a.revision) FROM doc_schema_revision a WHERE schema_id = ?`

	//get latest revision number from database
	row := db.QueryRow(sqlStr1, schemaInfo.ID)
	var tmpRev sql.NullInt64
	var revision int
	fetchErr := row.Scan(&tmpRev)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			revision = 0
		}

		return 0, fmt.Errorf("failed to fetch record from database: %s", fetchErr.Error())
	}

	if tmpRev.Valid {
		revision = int(tmpRev.Int64)
	} else {
		revision = 0
	}

	revision++ //increament by one for new revision
	oriRev := doc.Revision
	oriID := doc.ID
	oriName := doc.Name

	doc.Revision = revision
	doc.ID = schemaInfo.ID
	doc.Name = schemaInfo.Name

	xmlStr, xmlErr := doc.XML()

	doc.Revision = oriRev
	doc.ID = oriID
	doc.Name = oriName
	if xmlErr != nil {
		return 0, fmt.Errorf("failed to get XML definition: %s", xmlErr.Error())
	}

	_, insertErr := db.Exec(`INSERT INTO doc_schema_revision (schema_id,revision,xml_definition,remark) VALUES (?,?,?,?)`,
		schemaInfo.ID, revision, xmlStr, remark)

	if insertErr != nil {
		return 0, fmt.Errorf("failed to register new %s definition into database: %s",
			doc.Name, insertErr.Error())
	}

	return revision, nil
}

//SaveSchemaAsDraft save document schema as draft.
//	draft doc schema shall not affect production record
//NOTE: Existed draft record with same document name will be overwrite
//NOTE: ErrSchemaInfoNotFound error will return if document not register yet in doc_schema_info datatable
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

//SaveDraftToNewRevision convert draft into new revision
//Will return ErrDraftNotFound error if no draft available
func SaveDraftToNewRevision(db rdbmstool.DbHandlerProxy, name string) error {
	row := db.QueryRow(`SELECT COUNT(name) FROM doc_schema WHERE name = ? AND revision = -1`, name)
	var tmpInt int
	rowErr := row.Scan(&tmpInt)
	if rowErr != nil {
		if rowErr == sql.ErrNoRows {
			return ErrDraftNotFound{msg: fmt.Sprintf("no draft found for %s", name)}
		}

		return fmt.Errorf("failed to fetch record from database: %s", rowErr.Error())
	}

	if tmpInt == 0 {
		return ErrDraftNotFound{msg: fmt.Sprintf("no draft found for %s", name)}
	}

	info, infoErr := GetSchemaInfo(db, name)
	if infoErr != nil {
		return infoErr
	}

	_, updateErr := db.Exec(
		`UPDATE doc_schema SET revision = ? WHERE name = ? AND revision = -1`,
		info.LatestRevision+1, name)
	if updateErr != nil {
		return fmt.Errorf("failed to convert %s draft mode to release revision: %s",
			info.Name, updateErr.Error())
	}

	return nil
}

//UpdateDraftSchema update draft's xml definition
//Will return ErrDraftNotFound error if no draft in database
func UpdateDraftSchema(db rdbmstool.DbHandlerProxy, schema *gxschema.DxDoc) error {
	info, infoErr := GetSchemaByRevision(db, schema.Name, -1)
	if infoErr != nil {
		return fmt.Errorf("failed to fetch data from database: %s", infoErr.Error())
	}
	if info == nil {
		return ErrDraftNotFound{msg: fmt.Sprintf("no %s draft found in database", schema.Name)}
	}

	xmlStr, xmlErr := schema.XML()
	if xmlErr != nil {
		return fmt.Errorf("failed to generate XML string: %s", xmlErr.Error())
	}

	_, updateErr := db.Exec(
		`UPDATE doc_schema SET xml_definition = ? WHERE name = ? AND revision = -1`,
		xmlStr, schema.Name)
	if updateErr != nil {
		return fmt.Errorf("failed to update %s draft: %s", schema.Name, updateErr.Error())
	}

	return nil
}
