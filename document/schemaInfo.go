package document

import (
	"database/sql"
	"fmt"

	"github.com/guinso/rdbmstool"
)

//SchemaInfo summary of document schema
type SchemaInfo struct {
	Name           string
	LatestRevision int
	Description    string
	IsActive       bool
	HasDraft       bool
}

//GetSchemaInfo get SchemaInfo
func GetSchemaInfo(db rdbmstool.DbHandlerProxy, name string) (*SchemaInfo, error) {
	sqlStr := `SELECT a.name, a.description,  a.is_active, 
	MAX(b.revision), SUM(CASE  WHEN b.revision = -1 THEN 1 ELSE 0 END)
	FROM doc_schema_info a
	LEFT JOIN doc_schema b ON a.name = b.name
	WHERE a.name = ?
	GROUP BY a.name`

	row := db.QueryRow(sqlStr, name)
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
	SELECT a.name, a.description,  a.is_active, 
		MAX(b.revision), SUM(CASE  WHEN b.revision = -1 THEN 1 ELSE 0 END)
	FROM doc_schema_info a
	LEFT JOIN doc_schema b ON a.name = b.name
	GROUP BY a.name`

	rows, rowsErr := db.Query(sqlStr)
	if rowsErr != nil {
		return nil, fmt.Errorf("error encounter access database: %s", rowsErr.Error())
	}

	defer rows.Close()

	var tmpName, tmpDesc string
	var tmpLatestRev, tmpIsActive, tmpHasDraft int
	results := []SchemaInfo{}
	for rows.Next() {
		scanErr := rows.Scan(&tmpName, &tmpDesc, &tmpIsActive, &tmpLatestRev, &tmpHasDraft)
		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				break
			} else {
				return nil, fmt.Errorf("failed to fetch record from database: %s", scanErr.Error())
			}
		}

		results = append(results, SchemaInfo{
			Name:           tmpName,
			LatestRevision: tmpLatestRev,
			Description:    tmpDesc,
			IsActive:       tmpIsActive == 1,
			HasDraft:       tmpHasDraft == 1,
		})
	}

	return results, nil
}

//AddSchemaInfo register a new SchemaInfo into database,
//	if already exists will received ErrSchemaInfoAlreadyExists error
func AddSchemaInfo(db rdbmstool.DbHandlerProxy, name string, description string) error {
	sqlStr := `SELECT COUNT(name) FROM doc_schema_info WHERE name = ?`

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

	sqlInsert := `INSERT INTO doc_schema_info (name, description, is_active) VALUES (?,?,1)`
	_, dbErr := db.Exec(sqlInsert, name, description)
	if dbErr != nil {
		return fmt.Errorf("failed to create %s schemaInfo into database: %s", name, dbErr.Error())
	}

	return nil
}
