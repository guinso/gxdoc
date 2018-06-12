package initializeSequence

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/guinso/gxdoc/util"
	"github.com/guinso/gxschema"
)

func measureExecTime(startTime time.Time) {
	elapsed := time.Since(startTime)
	log.Println(fmt.Sprintf("exec time: %s", elapsed))
}

//routeProductionAPI handle dynamic HTTP user requset
func routeProductionAPI(w http.ResponseWriter, r *http.Request, trimURL string) {

	/***********************************************/
	//TODO: add your custom web API here:
	/**********************************************/

	//handle authentication web API
	//1. /login (POST)
	//2. /logout (POST)
	// if authentication.HandleHTTPRequest(util.GetDB(), w, r, trimURL) {
	// 	return
	// }

	//TODO: handle authorization web API
	//1. /role (GET + POST)
	// if authorization.HandleHTTPRequest(util.GetDB(), w, r, trimURL) {
	// 	return
	// }

	if strings.Compare(trimURL, "invoice") == 0 {
		rawInput := make(map[string]interface{})

		defer measureExecTime(time.Now())

		decodeErr := util.DecodeJSON(r, &rawInput)
		if decodeErr != nil {
			log.Println(decodeErr.Error())
			util.SendHTTPErrorResponse(w)
			return
		}

		defRaw := `
		<dxdoc name="invoice" revision="1">
			<dxstr name="invNo"></dxstr>
			<dxint name="totalQty" isOptional="true"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxdoc>`

		dxdoc, dxErr := gxschema.DecodeDxXML(defRaw)
		if dxErr != nil {
			log.Println(dxErr.Error())
			util.SendHTTPErrorResponse(w)
			return
		}

		// dxdoc := document.DxDoc{
		// 	Name:     "invoice",
		// 	Revision: 1,
		// 	Items: []document.DxItem{
		// 		document.DxStr{Name: "invNo", EnableLenLimit: true, LenLimit: 4},
		// 		document.DxInt{Name: "totalQty", IsOptional: true},
		// 		document.DxDecimal{Name: "price", Precision: 2},
		// 	},
		// }

		validateErr := dxdoc.ValidateData(rawInput)
		// elapsed := time.Since(startTime)
		// log.Println(fmt.Sprintf("exec time: %s", elapsed))

		if validateErr != nil {
			log.Println(validateErr.Error())
			util.SendHTTPResponse(w, -1, "invalid data format",
				fmt.Sprintf("{\"detail\": \"%s\"}", validateErr.Error()))
			return
		}

		util.SendHTTPResponse(w, 0, "ok", "{}")

	} else if strings.HasPrefix(trimURL, "meals") { //sample return JSON
		w.Header().Set("Content-Type", "application/json")  //MIME to application/json
		w.WriteHeader(http.StatusOK)                        //status code 200, OK
		w.Write([]byte("{ \"msg\": \"this is meal A \" }")) //body text
		return
	} else if strings.HasPrefix(trimURL, "img/") { //sample return virtual JPG file to client
		logicalFilePath := "./logic-files/"
		physicalFileName := "neon.jpg"

		// try read file
		data, err := ioutil.ReadFile(logicalFilePath + physicalFileName)
		if err != nil {
			// show error page if failed to read file
			handleErrorCode(500, "Unable to retrieve image file", w)
		} else {
			//w.Header().Set("Content-Type", "image/jpg") // #optional HTTP header info

			// uncomment if image file is meant to download instead of display on web browser
			// clientDisplayFileName = "customName.jpg"
			//w.header().Set("Content-Disposition", "attachment; filename=\"" + clientDisplayFileName + "\"")

			// write file (in binary format) direct into HTTP return content
			w.Write(data)
		}
	} else {
		// show error code 404 not found
		//(since the requested URL doesn't match any of it)
		handleErrorCode(404, "Path not found.", w)
	}
}
