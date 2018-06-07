package document

import (
	"testing"

	"github.com/guinso/gxdoc/testutil"
)

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
