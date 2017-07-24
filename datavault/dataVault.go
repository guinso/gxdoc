package datavault

import (
	"database/sql"
	"fmt"

	mysqlMeta "github.com/guinso/gxdoc/datavault/metareader/mysql"
	"github.com/guinso/gxdoc/datavault/record"

	///explicitly include GO mysql library
	_ "github.com/go-sql-driver/mysql"
)

///Gx DataVault
type DataVault struct {
	DbName    string
	DbAddress string
	Db        *sql.DB
}

//CreateDV to create DataVault instance
func CreateDV(address string, username string, password string,
	dbName string, port int) (*DataVault, error) {

	//TODO:  handle various database vendor
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

//GetHubs list all data vault's hub
func (dv *DataVault) GetHubs() []string {
	//TODO:  handle various database vendor
	return mysqlMeta.GetDbMetaTableName(dv.Db, dv.DbName, "hub_")
}

//GetSatalites list all data vault's satelites
func (dv *DataVault) GetSatalites() []string {
	//TODO:  handle various database vendor
	return mysqlMeta.GetDbMetaTableName(dv.Db, dv.DbName, "sat_")
}

//GetLinks list all data vault's links
func (dv *DataVault) GetLinks() []string {
	//TODO:  handle various database vendor
	return mysqlMeta.GetDbMetaTableName(dv.Db, dv.DbName, "link_")
}

//InsertRecord to insert new record into database
func (dv *DataVault) InsertRecord(dvInsertRecord *record.DvInsertRecord) error {
	sqls, sqlErr := dvInsertRecord.GenerateMultiSQL()

	if sqlErr != nil {
		return sqlErr
	}

	transaction, beginErr := dv.Db.Begin()
	if beginErr != nil {
		return beginErr
	}

	//TODO: test with various database vendor
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
