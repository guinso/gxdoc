package mysql

import (
	"database/sql"

	"fmt"

	"github.com/guinso/gxdoc/datavault/definition"
	"github.com/guinso/gxdoc/util"
)

//MetaReader implementation of MetaReader interface
type MetaReader struct {
	db     *sql.DB
	dbName string
}

//GetHubDefinition to get latest hub metainfo based on hub name
func (reader *MetaReader) GetHubDefinition(hubName string) (*definition.HubDefinition, error) {
	hubDef := definition.HubDefinition{}

	revision, revErr := reader.getDvEntityRevisionNumber("hub", hubName)
	if revErr != nil {
		return nil, revErr
	}

	//TODO: list all table columns

	return &hubDef, nil
}

//GetLinkDefinition to get latest link metainfo based on link name
func (reader *MetaReader) GetLinkDefinition(linkName string) (*definition.LinkDefinition, error) {
	return nil, nil
}

//GetSateliteDefinition to get latest satelite metainfo based on satelite name
func (reader *MetaReader) GetSateliteDefinition(satName string) (*definition.SateliteDefinition, error) {
	return nil, nil
}

//GetDbMetaTableName to get list of datatables' name which start with provided keyword
func (reader *MetaReader) GetDbMetaTableName(db *sql.DB, databaseName string, tableNamePrefix string) []string {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables"+
		" where table_schema=? AND table_name LIKE '"+tableNamePrefix+"%'", databaseName)

	if err != nil {
		return nil
	}

	var result []string
	for rows.Next() {
		var tmp string
		rows.Scan(&tmp)

		result = append(result, tmp)
	}

	return result
}

func (reader *MetaReader) getDvEntityRevisionNumber(entityName string, entityType string) (int, error) {
	//TODO: check table exists or not
	hubDbName := util.ToSnakeCase(entityName)
	tableName := fmt.Sprintf("%s_%s_rev", entityType, hubDbName)
	position := len(tableName) + 1

	sql := fmt.Sprintf("SELECT CONVERT(SUBSTRING(table_name, %d), SIGNED INTEGER) "+
		"AS rev FROM information_schema.tables "+
		"WHERE table_schema=? AND table_name LIKE '%s%%' ORDER BY rev DESC", position, tableName)
	revs, revErr := reader.db.Query(sql, reader.dbName)

	if revErr != nil {
		return -1, revErr
	}
	var revArr []int
	var tmp int
	for revs.Next() {
		revs.Scan(&tmp)

		revArr = append(revArr, tmp)
	}
	if len(revArr) == 0 {
		return -1, fmt.Errorf("entity %s not found in database", entityName)
	}

	return revArr[0], nil
}

func execSQL(sql string, transaction *sql.Tx) error {
	_, execErr := transaction.Exec(sql)
	if execErr != nil {
		return execErr
	}

	return nil
}
