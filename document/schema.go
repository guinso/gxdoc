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
	var tmpID string
	err := row.Scan(&xmlDef, &tmpRev, &tmpID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil //return NULL for record not found in database
		}

		return nil, fmt.Errorf("Failed to fetch doc_schema %s from DB: %s", name, err.Error())
	}

	//convert into DxDoc instance
	dxdoc, dxErr := gxschema.ParseSchemaFromXML(xmlDef)
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
	var tmpID string
	err := row.Scan(&xmlDef, &tmpID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil //return NULL for record not found
		}

		return nil, fmt.Errorf("failed to fetch doc_schema %s rev%d from database: %s",
			name, revision, err.Error())
	}

	//convert into DxDoc instance
	dxdoc, dxErr := gxschema.ParseSchemaFromXML(xmlDef)
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
func GetSchemaByID(db rdbmstool.DbHandlerProxy, id string) (*gxschema.DxDoc, error) {
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

		return nil, fmt.Errorf("Failed to fetch doc_schema ID %s from DB : %s", id, err.Error())
	}

	//convert into DxDoc instance
	dxdoc, dxErr := gxschema.ParseSchemaFromXML(xmlDef)
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
//NOTE: ErrSchemaInfoNotFound error will return if document not register yet in doc_schema datatable
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
//NOTE: if draft already exists, XML definition and remark will be overwriten
//NOTE: ErrSchemaInfoNotFound error will return if document not register in doc_schema datatable
func SaveSchemaAsDraft(db rdbmstool.DbHandlerProxy, schemaName string, doc *gxschema.DxDoc, remark string) error {
	xmlStr, xmlErr := doc.XML()
	if xmlErr != nil {
		return fmt.Errorf("failed convert doc schema into XML schema format: %s", xmlErr.Error())
	}

	//check SchemaInfo is registered
	schemaInfo, infoErr := GetSchemaInfo(db, schemaName)
	if infoErr != nil {
		return infoErr
	}
	if schemaInfo == nil {
		return ErrSchemaInfoNotFound{msg: schemaName + " not found in database"}
	}

	//check draft record is registered on database or not
	sqlStr1 := `SELECT COUNT(schema_id) FROM doc_schema_revision WHERE schema_id = ? AND revision = -1`
	row := db.QueryRow(sqlStr1, schemaInfo.ID)
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
		insertSQL := `INSERT INTO doc_schema_revision (schema_id, revision, xml_definition, remark) VALUES(?,-1,?,?)`
		_, dbErr := db.Exec(insertSQL, schemaInfo.ID, xmlStr, remark)
		if dbErr != nil {
			return fmt.Errorf("failed to save %s as draft into database: %s", schemaName, dbErr.Error())
		}
	} else {
		//update XML_dfinition column
		updateSQL := `UPDATE doc_schema_revision SET xml_definition = ?, remark = ? WHERE schema_id = ? AND revision = -1`
		_, dbErr := db.Exec(updateSQL, xmlStr, remark, schemaInfo.ID)
		if dbErr != nil {
			return fmt.Errorf("failed to save %s as draft into database: %s", schemaName, dbErr.Error())
		}
	}

	return nil
}

//SaveDraftToNewRevision convert draft into new revision
//Will return ErrDraftNotFound error if no draft available
func SaveDraftToNewRevision(db rdbmstool.DbHandlerProxy, schemaName string) error {
	//check SchemaInfo is registered
	schemaInfo, infoErr := GetSchemaInfo(db, schemaName)
	if infoErr != nil {
		return infoErr
	}
	if schemaInfo == nil {
		return ErrSchemaInfoNotFound{msg: schemaName + " not found in database"}
	}

	row := db.QueryRow(`SELECT COUNT(schema_id) FROM doc_schema_revision WHERE schema_id = ? AND revision = -1`, schemaInfo.ID)
	var tmpInt int
	rowErr := row.Scan(&tmpInt)
	if rowErr != nil {
		if rowErr == sql.ErrNoRows {
			return ErrDraftNotFound{msg: fmt.Sprintf("no draft found for %s", schemaName)}
		}

		return fmt.Errorf("failed to fetch record from database: %s", rowErr.Error())
	}

	if tmpInt == 0 {
		return ErrDraftNotFound{msg: fmt.Sprintf("no draft found for %s", schemaName)}
	}

	_, updateErr := db.Exec(
		`UPDATE doc_schema_revision SET revision = ? WHERE schema_id = ? AND revision = -1`,
		schemaInfo.LatestRevision+1, schemaInfo.ID)
	if updateErr != nil {
		return fmt.Errorf("failed to convert %s draft mode to release revision: %s",
			schemaInfo.Name, updateErr.Error())
	}

	return nil
}
