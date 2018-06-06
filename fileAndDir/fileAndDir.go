package fileAndDir

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/guinso/gxdoc/configuration"
	"github.com/guinso/gxdoc/util"
)

//InitStaticAndLogicDirectories initialize static and logic directories
func InitStaticAndLogicDirectories(config *configuration.ConfigInfo) error {
	//development static directory
	exists, err := util.IsDirectoryExists("./dev-static-files")
	if err != nil {
		fmt.Printf("Failed to check development static directory: %s", err.Error())
		return err
	}
	if !exists {
		if err = os.Mkdir("./dev-static-files", 0777); err != nil {
			fmt.Printf("Failed to create development static directory: %s", err.Error())
			return err
		}
	}

	//production logical directory
	exists, err = util.IsDirectoryExists(config.LogicDir)
	if err != nil {
		fmt.Printf("Failed to check logical directory: %s", err.Error())
		return err
	}
	if !exists {
		if err = os.Mkdir(config.LogicDir, 0777); err != nil {
			fmt.Printf("Failed to create logical directory: %s", err.Error())
			return err
		}
	}

	//production static directory
	exists, err = util.IsDirectoryExists(config.StaticDir)
	if err != nil {
		fmt.Printf("Failed to check static directory: %s", err.Error())
		return err
	}
	if !exists {
		if err = os.Mkdir(config.StaticDir, 0777); err != nil {
			fmt.Printf("Failed to create static directory: %s", err.Error())
			return err
		}

		//create basic index.html file
		htmlContent :=
			`<html><body>
				<h1>This is auto generated home page</h1>
				<p>To customize content, plase go to ./` + config.StaticDir + `</p>
			</body></html>`
		if err = ioutil.WriteFile(config.StaticDir+"/index.html", []byte(htmlContent), 0777); err != nil {
			return err
		}
	}

	return nil
}
