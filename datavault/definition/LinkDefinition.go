package definition

import (
	"errors"
	"fmt"

	"github.com/guinso/gxdoc/database"
	"github.com/guinso/gxdoc/util"
)

//LinkDefinition is schema to descibe link structure
type LinkDefinition struct {
	Name          string
	Revision      int
	HubReferences []HubReference
}

//GetHashKey is to generate data table equivalent hash key column name
func (linkDef *LinkDefinition) GetHashKey() string {
	return fmt.Sprintf("%s_hash_key", util.ToSnakeCase(linkDef.Name))
}

//GetDbTableName is to generate equivalent data table name
func (linkDef *LinkDefinition) GetDbTableName() string {
	return fmt.Sprintf("link_%s_rev%d", util.ToSnakeCase(linkDef.Name), linkDef.Revision)
}

// GenerateSQL is to generate SQL statement based on link definition
func (linkDef *LinkDefinition) GenerateSQL() (string, error) {
	if linkDef == nil || linkDef.HubReferences == nil || len(linkDef.HubReferences) < 2 {
		return "", errors.New("link definition must has atleast two hub reference")
	}

	tableDef := database.TableDefinition{
		Name:        linkDef.GetDbTableName(),
		PrimaryKey:  []string{linkDef.GetHashKey()},
		UniqueKeys:  []database.UniqueKeyDefinition{},
		ForiegnKeys: []database.ForeignKeyDefinition{},
		Columns: []database.ColumnDefinition{
			createHashKeyColumn(linkDef.Name),
			createLoadDateColumn(),
			createRecordSourceColumn()}}

	for _, hubRef := range linkDef.HubReferences {
		tableDef.Columns = append(tableDef.Columns, createHashKeyColumn(hubRef.HubName))

		tableDef.Indices = append(tableDef.Indices,
			database.IndexKeyDefinition{ColumnNames: []string{hubRef.GetHashKey()}})
		tableDef.ForiegnKeys = append(tableDef.ForiegnKeys,
			database.ForeignKeyDefinition{
				ColumnName:          hubRef.GetHashKey(),
				ReferenceTableName:  hubRef.GetDbTableName(),
				ReferenceColumnName: hubRef.GetHashKey()})
	}

	sql, err := database.GenerateTableSQL(&tableDef)
	if err != nil {
		return "", err
	}

	return sql, nil
}
