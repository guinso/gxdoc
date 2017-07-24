package definition

import (
	"fmt"

	"github.com/guinso/gxdoc/database"
	"github.com/guinso/gxdoc/util"
)

//HubDefinition is schema to descibe hub structure
type HubDefinition struct {
	Name         string
	BusinessKeys []string
	Revision     int
}

//GetHashKey is to generate data table equivalent hash key column name
func (hubDef *HubDefinition) GetHashKey() string {
	return fmt.Sprintf("%s_hash_key", util.ToSnakeCase(hubDef.Name))
}

//GetDbTableName is to generate equivalent data table name
func (hubDef *HubDefinition) GetDbTableName() string {
	return fmt.Sprintf("hub_%s_rev%d", util.ToSnakeCase(hubDef.Name), hubDef.Revision)
}

// GenerateSQL is to generate SQL statement based on hub definition
func (hubDef *HubDefinition) GenerateSQL() (string, error) {
	var sql string

	tableDef := database.TableDefinition{
		Name:        hubDef.GetDbTableName(),
		PrimaryKey:  []string{hubDef.GetHashKey()},
		UniqueKeys:  []database.UniqueKeyDefinition{},
		ForiegnKeys: []database.ForeignKeyDefinition{},
		Indices:     []database.IndexKeyDefinition{},
		Columns: []database.ColumnDefinition{
			createHashKeyColumn(hubDef.Name),
			createLoadDateColumn(),
			createRecordSourceColumn()}}

	if len(hubDef.BusinessKeys) > 0 {
		var uks []string

		for _, bk := range hubDef.BusinessKeys {
			tableDef.Columns = append(tableDef.Columns,
				database.ColumnDefinition{Name: util.ToSnakeCase(bk),
					DataType: database.CHAR, Length: 100, IsNullable: false})

			uks = append(uks, util.ToSnakeCase(bk))
		}

		tableDef.UniqueKeys = append(tableDef.UniqueKeys,
			database.UniqueKeyDefinition{ColumnNames: uks})
	}

	sql, err := database.GenerateTableSQL(&tableDef)

	if err != nil {
		return "", err
	}

	return sql, nil
}
