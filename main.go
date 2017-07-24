package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/guinso/gxdoc/database"
	"github.com/guinso/gxdoc/datavault"
	"github.com/guinso/gxdoc/datavault/definition"
	"github.com/guinso/gxdoc/datavault/record"
	"github.com/guinso/gxdoc/util"
)

//NOTE: work in progress (this is not actual use case)
func main() {
	fmt.Println("======== GX Backend Service ========")

	//TODO build query and begin Data Vault service
	fmt.Println("Starting Data Vault Engine.....")

	fmt.Println("Begin initialize database")
	dv, err := datavault.CreateDV("localhost", "root", "", "test", 3306)

	if err != nil {
		fmt.Println("Init DB failed")
		panic(err)
	}

	fmt.Println("DB initialize success")

	/*hubs := dv.GetHubs()
	fmt.Println("Hubs found: ")
	for _, item := range hubs {
		fmt.Println(item)
	}*/

	/*
		tableDef := database.TableDefinition{
			Name: "invoice",
			Columns: []database.ColumnDefinition{
				database.ColumnDefinition{"id", database.INTEGER, 11, false, 0},
				database.ColumnDefinition{"name", database.VARCHAR, 100, false, 0},
				database.ColumnDefinition{"job_id", database.VARCHAR, 100, false, 0}},
			PrimaryKey: []string{"id"},
			ForiegnKeys: []database.ForeignKeyDefinition{
				database.ForeignKeyDefinition{"job_id", "job", "id"}},
			UniqueKeys: []database.UniqueKeyDefinition{
				database.UniqueKeyDefinition{[]string{"name", "job_id"}}}}

		sql, sqlErr := database.GenerateTableSQL(&tableDef)
		if sqlErr != nil {
			panic(sqlErr)
		}

		fmt.Println(sql)
	*/

	/*
		hubDef := datavault.HubDefinition{Name: "taxInvoice", BusinessKeys: []string{"invoiceNo"}, Revision: 0}
		sql, err := datavault.GenerateHubSQL(&hubDef)
		if err != nil {
			panic(err)
		}
		fmt.Println(sql)
	*/

	/*
		linkDef := datavault.LinkDefinition{Name: "invoiceItem", Revision: 0,
			HubReferences: []datavault.HubReference{
				datavault.HubReference{HubName: "taxInvoice", Revision: 0},
				datavault.HubReference{HubName: "OrderItem", Revision: 0}}}
		sql, err := linkDef.GenerateSQL(&linkDef)
		if err != nil {
			panic(err)
		}
		fmt.Println(sql)
	*/

	/*
		satDef := definition.SateliteDefinition{
			Name:         "invoiceDetail",
			Revision:     0,
			HubReference: &definition.HubReference{HubName: "taxInvoice", Revision: 0},
			Attributes: []definition.SateliteAttributeDefinition{
				definition.SateliteAttributeDefinition{
					Name: "firstName", DataType: database.VARCHAR,
					Length: 100, IsNullable: false, DecimalPrecision: 0},
				definition.SateliteAttributeDefinition{
					Name: "lastName", DataType: database.VARCHAR,
					Length: 100, IsNullable: false, DecimalPrecision: 0},
				definition.SateliteAttributeDefinition{
					Name: "address", DataType: database.STRING,
					Length: 0, IsNullable: false, DecimalPrecision: 0}}}
		sql, err := satDef.GenerateSQL()
		if err != nil {
			panic(err)
		}
		fmt.Println(sql)
	*/

	loadDate := time.Date(2017, time.June, 20, 14, 10, 0, 0, &time.Location{})
	invoiceHashKey := util.MakeMD5("inv004")
	customerHashKey := util.MakeMD5("cus003")

	satCustomerInsert := record.SateliteInsertRecord{
		SateliteName:    "customer",
		Revision:        0,
		RecordSource:    "INVOICE",
		HubName:         "customer",
		HubHashKeyValue: customerHashKey,
		LoadDate:        loadDate,
		Attributes: []record.SateliteAttrInsertRecord{
			record.SateliteAttrInsertRecord{
				AttributeName: "firstName",
				Value:         "John",
				Meta: &definition.SateliteAttributeDefinition{
					Name:             "firstName",
					DataType:         database.CHAR,
					Length:           32,
					IsNullable:       false,
					DecimalPrecision: 0}},
			record.SateliteAttrInsertRecord{
				AttributeName: "lastName",
				Value:         "Doe",
				Meta: &definition.SateliteAttributeDefinition{
					Name:             "lastName",
					DataType:         database.CHAR,
					Length:           32,
					IsNullable:       false,
					DecimalPrecision: 0}}}}

	/*		gg, errrr := insertSat.GenerateSQL()
			if errrr != nil {
				panic(errrr)
			}
			fmt.Println(gg)*/

	hubInvoiceInsert := record.HubInsertRecord{
		HubName:      "invoice",
		HubRevision:  0,
		RecordSource: "INVOICE",
		LoadDate:     loadDate,
		HashKey:      invoiceHashKey,
		BusinessKeyVues: []record.HubBusinessKeyInsertRecord{
			record.HubBusinessKeyInsertRecord{
				BusinessKey:   "invoiceNumber",
				BusinessValue: "inv004"}}}

	hubCustomerInsert := record.HubInsertRecord{
		HubName:      "customer",
		HubRevision:  0,
		RecordSource: "INVOICE",
		LoadDate:     loadDate,
		HashKey:      customerHashKey,
		BusinessKeyVues: []record.HubBusinessKeyInsertRecord{
			record.HubBusinessKeyInsertRecord{
				BusinessKey:   "customerId",
				BusinessValue: "cus001"}}}

	/*hubSQL, hubErr := hubInsertRecord.GenerateSQL()
	if hubErr != nil {
		panic(hubErr)
	}
	fmt.Println(hubSQL)*/

	linkInsert := record.LinkInsertRecord{
		LinkName:     "invoiceItem",
		LinkRevision: 0,
		HashKey:      util.MakeMD5("1"),
		RecordSource: "INVOICE",
		LoadDate:     loadDate,
		ReferenceHashKey: []record.LinkReferenceInsertRecord{
			record.LinkReferenceInsertRecord{
				HubName:      "invoice",
				HashKeyValue: invoiceHashKey},
			record.LinkReferenceInsertRecord{
				HubName:      "customer",
				HashKeyValue: customerHashKey}}}
	/*	linkSQL, linkErr := linkInsertRecord.GenerateSQL()
		if linkErr != nil {
			panic(linkErr)
		}
		fmt.Println(linkSQL)*/

	dvInsert := record.DvInsertRecord{
		LoadDate:  loadDate,
		Hubs:      []record.HubInsertRecord{hubCustomerInsert, hubInvoiceInsert},
		Links:     []record.LinkInsertRecord{linkInsert},
		Satelites: []record.SateliteInsertRecord{satCustomerInsert}}

	/*dvSQL, dvErr := dvInsert.GenerateSQL()
	if dvErr != nil {
		panic(dvErr)
	}
	fmt.Println(dvSQL)*/

	sqls, sqlErr := dvInsert.GenerateMultiSQL()
	if sqlErr != nil {
		panic(sqlErr)
	}

	for _, sql := range sqls {
		fmt.Println(sql)
	}

	insertErr := dv.InsertRecord(&dvInsert)
	if insertErr != nil {
		panic(insertErr)
	}
	fmt.Println("Insert Datavault record success.")

	/*
		// handler all request start from "/"
		http.HandleFunc("/", handler)

		// start HTTP server in socket 7777
		err := http.ListenAndServe(":7777", nil)

		// start HTTPS server (default socket 443)
		//x err := http.ListenAndServeTLS(":7777", "chetsiang2.crt", "chetsiang2.key", nil)

		if err != nil {
			errMsg := err.Error()
			fmt.Println(errMsg)
		}
	*/
}

// Handle HTTP request to either static file server or REST server (URL start with "api/")
func handler(w http.ResponseWriter, r *http.Request) {
	//remove first "/" character
	urlPath := r.URL.Path[1:]

	//if start with "api/" direct to REST handler
	if strings.HasPrefix(urlPath, "api/") {
		//trim prefix "api/"
		trimmedURL := urlPath[4:]
		/*
			//trim suffix "/"
			if strings.HasSuffix(trimmedURL, "/") {
				trimmedURL = trimmedURL[0:(len(trimmedURL) - 1)]
			}
		*/

		routePath(w, r, trimmedURL)
	} else {
		log.Print("Entering static file handler: " + urlPath)

		// define your static file directory
		staticFilePath := "./static-files/"

		//other wise, let read a file path and display to client
		http.ServeFile(w, r, staticFilePath+urlPath)
	}
}

//handle dynamic HTTP user requset
func routePath(w http.ResponseWriter, r *http.Request, trimURL string) {
	// find and match trimmed URL to respective REST request
	//      trimmed URL
	//      request method
	//      input parameter(s)

	if strings.HasPrefix(trimURL, "login") {
		// example URL: localhost:7777/api/login
		// TODO: handle login request
		fmt.Fprint(w, "Request login")
	} else if strings.HasPrefix(trimURL, "logout") {
		// example URL: localhost:7777/api/logout
		// TODO: handle logout request
		fmt.Fprint(w, "Request logout")
	} else if strings.HasPrefix(trimURL, "meals") {
		// example URL: localhost:7777/api/meals
		// show list of meals
		w.Header().Set("Content-Type", "application/json")                                        //MIME to application/json
		w.WriteHeader(http.StatusOK)                                                              //status code 200, OK
		w.Write([]byte(`{ 
			"msg": "this is meal A",
			"id": "A123",
			"name": "Steak"
		}`)) //body text
	} else {
		// show error code 404 not found
		util.HandleErrorCode(404, "Path not found.", w)
	}
}
