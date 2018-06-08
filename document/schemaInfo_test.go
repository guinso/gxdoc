package document

import (
	"strings"
	"testing"

	"github.com/guinso/gxdoc/testutil"
)

func TestGetSchemaInfo(t *testing.T) {
	db, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	tmp, tmpErr := GetSchemaInfo(db, "invoice")
	if tmpErr != nil {
		t.Error(tmpErr)
	}
	if tmp == nil {
		t.Errorf("expect invoice is registed on database")
	}

	tmp, tmpErr = GetSchemaInfo(db, "invoice123")
	if tmpErr != nil {
		t.Error(tmpErr)
	}
	if tmp != nil {
		t.Errorf("expect invoice is not registed on database")
	}
}

func TestGetSchemaInfoByID(t *testing.T) {
	db, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	tmp, tmpErr := GetSchemaInfoByID(db, 1)
	if tmpErr != nil {
		t.Error(tmpErr)
	}
	if tmp == nil {
		t.Errorf("expect invoice is registed on database")
	}

	tmp, tmpErr = GetSchemaInfoByID(db, 3)
	if tmpErr != nil {
		t.Error(tmpErr)
	}
	if tmp != nil {
		t.Errorf("expect no record (ID 3) is registed on database")
	}
}

func TestGetAllSchemaInfo(t *testing.T) {
	db, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	items, itemsErr := GetAllSchemaInfo(db)
	if itemsErr != nil {
		t.Error(itemsErr)
	}

	if len(items) != 2 {
		t.Errorf("expect database has 2 document schema but get %d instead", len(items))
		return
	}

	if items[1].HasDraft != true {
		t.Errorf("expect items[1] has draft mode")
	}

	if items[0].HasDraft != false {
		t.Errorf("expect items[0] has no draft mode")
	}
}

func TestUpdateSchemaInfo(t *testing.T) {
	db, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	trx, trxErr := db.Begin()
	if trxErr != nil {
		t.Fatal(trxErr)
		return
	}

	defer trx.Rollback()

	invInfo := SchemaInfo{
		Name:        "invoice 123",
		ID:          1,
		Description: "blah...",
		IsActive:    false,
	}

	updateErr := UpdateSchemaInfo(trx, &invInfo)
	if updateErr != nil {
		t.Error(updateErr)
		return
	}

	newInvInfo, newErr := GetSchemaInfo(trx, "invoice 123")
	if newErr != nil {
		t.Error(newErr)
		return
	}

	if newInvInfo == nil {
		t.Errorf("newly update invoice name (invoice 123) should be found in database")
		return
	}

	if newInvInfo.ID != 1 {
		t.Errorf("expect invoice ID is 1 but get %d", newInvInfo.ID)
	}

	if newInvInfo.IsActive == true {
		t.Errorf("expect invoice is not active")
	}

	if strings.Compare(newInvInfo.Description, "blah...") != 0 {
		t.Errorf("expect invoice description is 'blah...' but get '%s'", newInvInfo.Description)
	}
}

func TestAddSchemaInfo(t *testing.T) {
	db, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	trx, trxErr := db.Begin()
	if trxErr != nil {
		t.Fatal(trxErr)
		return
	}

	defer trx.Rollback()

	addErr := AddSchemaInfo(trx, "po", "purchase order")
	if addErr != nil {
		t.Error(addErr)
		return
	}

	po, poErr := GetSchemaInfo(trx, "po")
	if poErr != nil {
		t.Error(poErr)
		return
	}

	if po.ID != 3 {
		t.Errorf("expect PO ID is 3 but get %d", po.ID)
	}
	if strings.Compare(po.Description, "purchase order") != 0 {
		t.Errorf("expect PO description is 'purchase order' but get '%s'", po.Description)
	}
}
