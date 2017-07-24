package datavault

import (
	"database/sql"
	"fmt"
	"gxdoc/datavault/record"

	///explicitly include GO mysql library
	_ "github.com/go-sql-driver/mysql"
)

///Gx DataVault
type DataVault struct {
	DbName    string
	DbAddress string
	Db        *sql.DB
}

///Initialize DataVault instance
func CreateDV(address string, username string, password string,
	dbName string, port int) (*DataVault, error) {

	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, address, port, dbName))

	if err != nil {
		return nil, err
	}

	pingErr := db.Ping() //check connection is valid or not

	if pingErr != nil {
		return nil, pingErr
	}

	dv := DataVault{dbName, address, db}

	return &dv, nil
}

///list all data vault's hub
func (dv *DataVault) GetHubs() []string {
	return dv.getDbMetaTableName("hub_")
}

///list all data vault's satelites
func (dv *DataVault) GetSatalites() []string {
	return dv.getDbMetaTableName("sat_")
}

///list all data vault's links
func (dv *DataVault) GetLinks() []string {
	return dv.getDbMetaTableName("link_")
}

func (dv *DataVault) getDbMetaTableName(tablePrefix string) []string {
	rows, err := dv.Db.Query("SELECT table_name FROM information_schema.tables"+
		" where table_schema=? AND table_name LIKE '"+tablePrefix+"%'", dv.DbName)

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

//InsertRecord is use to insert new record into database
func (dv *DataVault) InsertRecord(dvInsertRecord *record.DvInsertRecord) error {
	sqls, sqlErr := dvInsertRecord.GenerateMultiSQL()

	if sqlErr != nil {
		return sqlErr
	}

	transaction, beginErr := dv.Db.Begin()
	if beginErr != nil {
		return beginErr
	}

	for _, sql := range sqls {
		execErr := dv.execSQL(sql, transaction)
		if execErr != nil {
			transaction.Rollback()
			return execErr
		}
	}

	commitErr := transaction.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}

func (dv *DataVault) execSQL(sql string, transaction *sql.Tx) error {
	_, execErr := transaction.Exec(sql)
	if execErr != nil {
		return execErr
	}

	return nil
}
