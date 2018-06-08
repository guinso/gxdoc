package document

import (
	"reflect"
	"strings"
	"testing"

	"github.com/guinso/gxdoc/testutil"
	"github.com/guinso/gxschema"
)

func TestGetSchema(t *testing.T) {
	testDb, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	invoice, dxErr := GetSchema(testDb, "invoice")
	if dxErr != nil {
		t.Fatal(dxErr)
		return
	}

	if invoice == nil {
		t.Fatal("expect invoice schema definition record is within database")
	}

	if invoice.Revision != 2 {
		t.Errorf("expect latest revision is 2 but get %d instead", invoice.Revision)
	}
	if strings.Compare(invoice.Name, "invoice") != 0 {
		t.Errorf("expect document name is 'invoice' but get '%s'", invoice.Name)
	}
	if invoice.ID != 1 {
		t.Errorf("expect ID is 1 but get %d", invoice.ID)
	}

	dxstr, strOK := invoice.Items[0].(*gxschema.DxStr)
	if !strOK {
		t.Errorf("expect first item is DxStr but get %s", reflect.TypeOf(invoice.Items[0]))
	} else if strings.Compare(dxstr.Name, "invNo") != 0 {
		t.Errorf("expect item[0] name id 'invNo' but get '%s'", dxstr.Name)
	}

	///
	xxx, xxxErr := GetSchema(testDb, "pr")
	if xxxErr != nil {
		t.Error(xxxErr)
	}
	if xxx == nil {
		t.Errorf("PR should found in doc_schema data table")
	}
}

func TestGetSchema2(t *testing.T) {
	testDb, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	invoice, dxErr := GetSchema(testDb, "invoice123")
	if dxErr != nil {
		t.Fatal(dxErr)
		return
	}

	if invoice != nil {
		t.Fatal("expect no schema definition record found in database")
	}
}

func TestGetSchemaByID(t *testing.T) {
	testDb, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	invoice, dxErr := GetSchemaByID(testDb, 1)
	if dxErr != nil {
		t.Fatal(dxErr)
		return
	}

	if invoice == nil {
		t.Fatal("expect invoice schema definition record is within database")
	}

	if invoice.Revision != 2 {
		t.Errorf("expect latest revision is 2 but get %d instead", invoice.Revision)
	}
	if strings.Compare(invoice.Name, "invoice") != 0 {
		t.Errorf("expect document name is 'invoice' but get '%s'", invoice.Name)
	}
	if invoice.ID != 1 {
		t.Errorf("expect ID is 1 but get %d", invoice.ID)
	}

	dxstr, strOK := invoice.Items[0].(*gxschema.DxStr)
	if !strOK {
		t.Errorf("expect first item is DxStr but get %s", reflect.TypeOf(invoice.Items[0]))
	} else if strings.Compare(dxstr.Name, "invNo") != 0 {
		t.Errorf("expect item[0] name id 'invNo' but get '%s'", dxstr.Name)
	}
}

func TestGetSchemaByRevision(t *testing.T) {
	testDb, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	invoice, dxErr := GetSchemaByRevision(testDb, "invoice", 1)
	if dxErr != nil {
		t.Fatal(dxErr)
		return
	}

	if invoice.Revision != 1 {
		t.Errorf("expect latest revision is 1 but get %d instead", invoice.Revision)
	}
	if invoice.ID != 1 {
		t.Errorf("expect ID is 1 but get %d", invoice.ID)
	}
}

func TestGetSchemaByRevision2(t *testing.T) {
	testDb, dbErr := testutil.GetTestDB()
	if dbErr != nil {
		t.Fatal(dbErr)
		return
	}

	invoice, dxErr := GetSchemaByRevision(testDb, "invoice123", 3)
	if dxErr != nil {
		t.Fatal(dxErr)
		return
	}

	if invoice != nil {
		t.Fatal("expect no schema definition record found in database")
	}
}

func TestAddSchema(t *testing.T) {
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

	doc := gxschema.DxDoc{
		Name:     "invoice",
		Revision: 0,
		Items: []gxschema.DxItem{
			gxschema.DxInt{Name: "qty123"},
		},
	}

	latestRev, addErr := AddSchema(trx, "invoice", &doc, "sample 1")
	if addErr != nil {
		t.Error(addErr)
		return
	}

	invoice, docErr := GetSchema(trx, "invoice")
	if docErr != nil {
		t.Error(docErr)
		return
	}

	if invoice.Revision != 3 {
		t.Errorf("expect latest revision for invoice is 3 but get %d", latestRev)
	}

	if len(invoice.Items) != 1 {
		t.Errorf("expect latest invoice only have one item definition but get %d", len(invoice.Items))
	} else {
		dxint, intOK := invoice.Items[0].(*gxschema.DxInt)
		if !intOK {
			t.Errorf("expect invoice.items[0] is DxInt but get %s", reflect.TypeOf(invoice.Items[0]))
			return
		}

		if strings.Compare(dxint.Name, "qty123") != 0 {
			t.Errorf("expect invoice.items[0] name is 'qty123' but get %s", dxint.Name)
		}
	}
}

func TestSaveSchemaAsDraft(t *testing.T) {
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

	doc := gxschema.DxDoc{
		Name:     "invoice",
		Revision: 0,
		Items: []gxschema.DxItem{
			gxschema.DxInt{Name: "qty123"},
		},
	}

	draftErr := SaveSchemaAsDraft(trx, &doc, "try save as new record")
	if draftErr != nil {
		t.Error(draftErr)
		return
	}

	draftInvoice, invErr := GetDraftSchema(trx, "invoice")
	if invErr != nil {
		t.Fatal(invErr)
		return
	}
	if draftInvoice == nil {
		t.Errorf("expect database has invoice draft record")
	}

	doc.Items = append(doc.Items, &gxschema.DxStr{Name: "description"})
	draftErr = SaveSchemaAsDraft(trx, &doc, "try save as an update")
	if draftErr != nil {
		t.Error(draftErr)
	}

	draftInvoice, invErr = GetDraftSchema(trx, "invoice")
	if invErr != nil {
		t.Fatal(invErr)
		return
	}
	if draftInvoice == nil {
		t.Errorf("expect database has invoice draft record after update record")
	}

	if len(draftInvoice.Items) != 2 {
		t.Errorf("expect latest draft version of invoice has 2 item definition but get %d",
			len(draftInvoice.Items))
	}
}

func TestGetDraftSchema(t *testing.T) {
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

	doc := gxschema.DxDoc{
		Name:     "invoice",
		Revision: 0,
		Items: []gxschema.DxItem{
			gxschema.DxInt{Name: "qty123"},
		},
	}

	draftDoc, docErr := GetDraftSchema(trx, "invoice")
	if docErr != nil {
		t.Error(docErr)
		return
	}
	if draftDoc != nil {
		t.Errorf("expect no draft document schema found in database")
		return
	}

	draftErr := SaveSchemaAsDraft(trx, &doc, "add draft record and try retrieve from database again")
	if draftErr != nil {
		t.Error(draftErr)
		return
	}
	draftDoc, docErr = GetDraftSchema(trx, "invoice")
	if docErr != nil {
		t.Error(docErr)
		return
	}
	if draftDoc == nil {
		t.Errorf("expect there is draft document schema in database")
		return
	}
}

func TestSaveDraftToNewRevision(t *testing.T) {
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

	darftErr := SaveDraftToNewRevision(trx, "pr")
	if darftErr != nil {
		t.Error(darftErr)
		return
	}

	prInfo, infoErr := GetSchemaByRevision(trx, "pr", -1)
	if infoErr != nil {
		t.Error(infoErr)
		return
	}
	if prInfo != nil {
		t.Errorf("exepct pr has no draft mode anyomre")
	}
}

func TestUpdateDraftSchema(t *testing.T) {
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

	dxdoc := gxschema.DxDoc{
		Name:     "pr",
		Revision: -1,
		Items: []gxschema.DxItem{
			gxschema.DxInt{Name: "qty"},
			gxschema.DxStr{Name: "pr number", EnableLenLimit: true, LenLimit: 6},
			gxschema.DxBool{Name: "mandatory"},
		},
	}

	updatErr := UpdateDraftSchema(trx, &dxdoc)
	if updatErr != nil {
		t.Error(updatErr)
		return
	}

	pr, prErr := GetDraftSchema(trx, "pr")
	if prErr != nil {
		t.Error(prErr)
		return
	}
	if len(pr.Items) != 3 {
		t.Errorf("expect PR draft has 3 definition but get %d instead", len(pr.Items))
	}
}
