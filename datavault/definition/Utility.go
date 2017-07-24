package definition

import (
	"gxdoc/database"
	"gxdoc/util"
)

const (
	//LOAD_DATE is data vault standard table column name
	LOAD_DATE = "load_date"
	//END_DATE is data vault standard table column name
	END_DATE = "end_date"
	//RECORD_SOURCE is data vault standard table column name
	RECORD_SOURCE = "record_source"
)

func createHashKeyColumn(name string) database.ColumnDefinition {
	return database.ColumnDefinition{
		Name:     util.ToSnakeCase(name) + "_hash_key",
		DataType: database.CHAR, Length: 32, IsNullable: false}
}

func createEndDateColumn() database.ColumnDefinition {
	return database.ColumnDefinition{Name: END_DATE,
		DataType: database.DATE, Length: 0, IsNullable: true}
}

func createLoadDateColumn() database.ColumnDefinition {
	return database.ColumnDefinition{Name: LOAD_DATE,
		DataType: database.DATE, Length: 0, IsNullable: false}
}

func createRecordSourceColumn() database.ColumnDefinition {
	return database.ColumnDefinition{Name: RECORD_SOURCE,
		DataType: database.CHAR, Length: 100, IsNullable: false}
}

func createIndexKey(colName string) database.IndexKeyDefinition {
	return database.IndexKeyDefinition{ColumnNames: []string{colName}}
}
