package database

import "database/sql"

//IsDbTableExists check the data table is found within selected database
func IsDbTableExists(tableName string, databaseName string, db *sql.DB) (bool, error) {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables"+
		" where table_schema=? AND table_name=?", databaseName, tableName)

	if err != nil {
		return false, err
	}

	var result int
	for rows.Next() {
		result++
	}

	return result > 0, nil
}

//SearchDbTableName search for matching data table name based on search keywords
func SearchDbTableName(searchQuery string, databaseName string, db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables "+
		"where table_schema=? AND table_name LIKE ?", databaseName, searchQuery)

	if err != nil {
		return nil, err
	}

	var result []string
	var tableName string
	for rows.Next() {
		scanErr := rows.Scan(&tableName)

		if scanErr != nil {
			return nil, scanErr
		}

		result = append(result, tableName)
	}

	return result, nil
}
