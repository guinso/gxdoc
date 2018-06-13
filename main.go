package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/guinso/gxdoc/bootSequence"
	"github.com/guinso/gxdoc/configuration"
	"github.com/guinso/gxdoc/util"
)

func main() {
	fmt.Print("loading configuration file...")
	configErr := configuration.LoadINIConfigFile()
	if configErr != nil {
		fmt.Println("\t\t\t[FAILED]")
		panic(configErr)
	}
	fmt.Println("\t\t\t[OK]")

	fmt.Print("try connect to MySQL database...")
	db, dbErr := checkDbConnection(configuration.GetConfig())
	if dbErr != nil {
		fmt.Println("\t\t[FAILED]")
		panic(dbErr)
	}
	util.SetDB(db)
	fmt.Println("\t\t[OK]")

	fmt.Print("creating directories...")
	if dirErr := bootSequence.InitStaticAndLogicDirectories(configuration.GetConfig()); dirErr != nil {
		fmt.Println("\t\t\t\t[FAILED]")
		panic(dirErr)
	}
	fmt.Println("\t\t\t\t[OK]")

	fmt.Println(fmt.Sprintf("starting web server on port %d", configuration.GetConfig().PortNumber))
	if webErr := startWebServer(); webErr != nil {
		panic(webErr)
	}
}

//checkDbConnection try connect to MySQL database ans check able to connect or not
func checkDbConnection(config *configuration.ConfigInfo) (*sql.DB, error) {
	//TODO:  handle various database vendor
	dbx, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		config.DbUsername,
		config.DbPassword,
		config.DbAddress,
		config.DbPort,
		config.DbName))

	if err != nil {
		return nil, err
	}

	//check connection is valid or not
	if pingErr := dbx.Ping(); pingErr != nil {
		return nil, pingErr
	}

	return dbx, nil
}

func startWebServer() error {
	config := configuration.GetConfig()

	bootSequence.SetConfig(config.StaticDir, config.DevEnable, config.DevStartURL, "static-files-dev")

	http.HandleFunc("/", bootSequence.HandleRouting)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.PortNumber), nil)
}
