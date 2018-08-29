package document

import (
	"database/sql"
	"fmt"

	"github.com/guinso/stringtool"

	"github.com/guinso/rdbmstool"
)

//SchemaInfo summary of document schema
type SchemaInfo struct {
	ID             string
	Name           string
	LatestRevision int
	Description    string
	IsActive       bool
	HasDraft       bool
}

//GetSchemaInfo get SchemaInfo
func GetSchemaInfo(db rdbmstool.DbHandlerProxy, name string) (*SchemaInfo, error) {
	sqlStr := `SELECT a.id, a.name, a.description,  a.is_active, 
	MAX(b.revision), SUM(CASE  WHEN b.revision = -1 THEN 1 ELSE 0 END)
	FROM doc_schema a
	LEFT JOIN doc_schema_revision b ON a.id = b.schema_id
	WHERE a.name = ?
	GROUP BY a.name`

	row := db.QueryRow(sqlStr, name)
	var tmpID, tmpName, tmpDesc string
	var tmpIsActive int
	var tmpLatestRev, tmpHasDraft sql.NullInt64
	scanErr := row.Scan(&tmpID, &tmpName, &tmpDesc, &tmpIsActive, &tmpLatestRev, &tmpHasDraft)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to fetch record from database: %s", scanErr.Error())
	}

	var finalRev int
	if tmpLatestRev.Valid {
		finalRev = int(tmpLatestRev.Int64)
	} else {
		finalRev = 0
	}

	var finalHasDraft bool
	if tmpHasDraft.Valid {
		finalHasDraft = tmpHasDraft.Int64 == 1
	} else {
		finalHasDraft = false
	}

	return &SchemaInfo{
		ID:             tmpID,
		Name:           tmpName,
		LatestRevision: finalRev,
		Description:    tmpDesc,
		IsActive:       tmpIsActive == 1,
		HasDraft:       finalHasDraft,
	}, nil
}

//GetSchemaInfoByID get schema info from database by ID
func GetSchemaInfoByID(db rdbmstool.DbHandlerProxy, IDD string) (*SchemaInfo, error) {
	sqlStr := `SELECT a.name, a.description,  a.is_active, 
	MAX(b.revision), SUM(CASE  WHEN b.revision = -1 THEN 1 ELSE 0 END)
	FROM doc_schema a
	LEFT JOIN doc_schema_revision b ON a.id = b.schema_id
	WHERE a.id = ?
	GROUP BY a.name`

	row := db.QueryRow(sqlStr, IDD)
	var tmpName, tmpDesc string
	var tmpLatestRev, tmpIsActive, tmpHasDraft int
	scanErr := row.Scan(&tmpName, &tmpDesc, &tmpIsActive, &tmpLatestRev, &tmpHasDraft)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to fetch record from database: %s", scanErr.Error())
	}

	return &SchemaInfo{
		ID:             IDD,
		Name:           tmpName,
		LatestRevision: tmpLatestRev,
		Description:    tmpDesc,
		IsActive:       tmpIsActive == 1,
		HasDraft:       tmpHasDraft == 1,
	}, nil
}

//GetAllSchemaInfo get all available document schema summary
func GetAllSchemaInfo(db rdbmstool.DbHandlerProxy) ([]SchemaInfo, error) {
	sqlStr := `
	SELECT a.id, a.name, a.description,  a.is_active, 
		MAX(b.revision), SUM(CASE  WHEN b.revision = -1 THEN 1 ELSE 0 END)
	FROM doc_schema a
	LEFT JOIN doc_schema_revision b ON a.id = b.schema_id
	GROUP BY a.id`

	rows, rowsErr := db.Query(sqlStr)
	if rowsErr != nil {
		return nil, fmt.Errorf("error encounter access database: %s", rowsErr.Error())
	}

	defer rows.Close()

	var tmpID, tmpName, tmpDesc string
	var tmpLatestRev, tmpIsActive, tmpHasDraft int
	results := []SchemaInfo{}
	for rows.Next() {
		scanErr := rows.Scan(&tmpID, &tmpName, &tmpDesc, &tmpIsActive, &tmpLatestRev, &tmpHasDraft)
		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				break
			} else {
				return nil, fmt.Errorf("failed to fetch record from database: %s", scanErr.Error())
			}
		}

		results = append(results, SchemaInfo{
			ID:             tmpID,
			Name:           tmpName,
			LatestRevision: tmpLatestRev,
			Description:    tmpDesc,
			IsActive:       tmpIsActive == 1,
			HasDraft:       tmpHasDraft == 1,
		})
	}

	return results, nil
}

//UpdateSchemaInfo update schema info description and isActive attributes
//Will return ErrSchemaInfoNotFound if specified schema info not registered yet on database
func UpdateSchemaInfo(db rdbmstool.DbHandlerProxy, docInfo *SchemaInfo) error {
	row := db.QueryRow(`SELECT COUNT(id) FROM doc_schema WHERE id = ?`, docInfo.ID)
	var tmpInt int
	err := row.Scan(&tmpInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrSchemaInfoNotFound{msg: fmt.Sprintf("%s not found in database", docInfo.Name)}
		}

		return err
	}

	if tmpInt == 0 {
		return ErrSchemaInfoNotFound{msg: fmt.Sprintf("%s not found in database", docInfo.Name)}
	}

	_, updateErr := db.Exec(
		`UPDATE doc_schema SET name = ?, description = ?, is_active = ? WHERE id = ?`,
		docInfo.Name, docInfo.Description, docInfo.IsActive, docInfo.ID)

	if updateErr != nil {
		return fmt.Errorf("failed to update schema info's description: %s", updateErr.Error())
	}

	return nil
}

//AddSchemaInfo register a new SchemaInfo into database,
//if already exists will received ErrSchemaInfoAlreadyExists error
func AddSchemaInfo(db rdbmstool.DbHandlerProxy, name string, description string) error {
	sqlStr := `SELECT COUNT(name) FROM doc_schema WHERE name = ?`

	row := db.QueryRow(sqlStr, name)
	var tmpInt int
	scanErr := row.Scan(&tmpInt)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			tmpInt = 0
		} else {
			return scanErr
		}
	}

	if tmpInt > 0 {
		return ErrSchemaInfoAlreadyExists{msg: name + " already exists"}
	}

	//get next ID
	// currentID, currErr := util.GetDBColumnMaxInt(db, "doc_schema", "id")
	// if currErr != nil {
	// 	return fmt.Errorf("failed to acquire latest record ID: %s", currErr.Error())
	// }
	// currentID++
	newID, idErr := stringtool.GenerateRandomUUID()
	if idErr != nil {
		return fmt.Errorf("failed to generate ID for new schema %s", name)
	}

	sqlInsert := `INSERT INTO doc_schema (id, name, description, is_active) VALUES (?,?,?,1)`
	_, dbErr := db.Exec(sqlInsert, newID, name, description)
	if dbErr != nil {
		return fmt.Errorf("failed to create %s schemaInfo into database: %s", name, dbErr.Error())
	}

	return nil
}

//JSON export to JSON string
func (info *SchemaInfo) JSON() string {
	return fmt.Sprintf(
		`{"name": "%s","latestRev": %d,"desc":`+
			` "%s","isActive": %t,"hasDraft": %t}`,
		info.Name, info.LatestRevision,
		info.Description, info.IsActive, info.HasDraft)
}
