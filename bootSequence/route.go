package bootSequence

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/guinso/gxdoc/routing"
	"github.com/guinso/gxdoc/util"
	"github.com/guinso/gxschema"
)

var enableDevMode = false
var devStartURL string

var staticPath string
var devStaticPath string

//SetConfig set routing configuration
func SetConfig(staticDirectory string, enableDev bool, developmentStartURL string,
	developmentStaticPath string) {

	staticPath = staticDirectory

	enableDevMode = enableDev
	devStartURL = developmentStartURL
	devStaticPath = developmentStaticPath
}

// HandleRouting HTTP request to either static file server or REST server (URL start with "api/")
func HandleRouting(w http.ResponseWriter, r *http.Request) {

	var urlPath = r.URL.Path

	//remove first "/" character
	if strings.HasPrefix(urlPath, "/") {
		urlPath = r.URL.Path[1:]
	}
	log.Println(fmt.Sprintf("Serving URL[%s]: %s", r.Method, r.URL.Path))

	//proceed to development if path matched
	if enableDevMode == true && (strings.HasPrefix(urlPath, devStartURL+"/") ||
		strings.Compare(urlPath, devStartURL) == 0) {

		urlPath = urlPath[len(devStartURL):]
		if strings.HasPrefix(urlPath, "/") {
			urlPath = urlPath[1:]
		}

		//if start with "api/" direct to REST handler
		if strings.HasPrefix(urlPath, "api/") || strings.Compare(urlPath, "api") == 0 {
			//log.Println("route handle by development API")
			if strings.HasPrefix(urlPath, "api/") {
				routeDevelopmentAPI(w, r, urlPath[4:]) //trim: 'api/'
			} else {
				routeDevelopmentAPI(w, r, "")
			}
		} else {
			//log.Println("route handle by development static file")
			//other wise, lets read a file path and display to client
			http.ServeFile(w, r, "./"+devStaticPath+"/"+urlPath)
		}

	} else {
		//if start with "api/" direct to REST handler
		if strings.HasPrefix(urlPath, "api/") || strings.Compare(urlPath, "api") == 0 {
			//log.Println("route handle by production API")
			if strings.HasPrefix(urlPath, "api/") {
				routeProductionAPI(w, r, urlPath[4:]) //trim: "api/"
			} else {
				routeProductionAPI(w, r, "")
			}

		} else {
			//log.Println("route handle by production stati file")
			//other wise, lets read a file path and display to client
			http.ServeFile(w, r, "./"+staticPath+"/"+urlPath)
		}
	}
}

//routeDevelopmentAPI handle HTTP REST API (development)
func routeDevelopmentAPI(w http.ResponseWriter, r *http.Request, url string) {
	//TODO: how to check requstor ID and determine her authority?

	done, docErr := HandleDocSchemaHTTP(url, w, r)
	if done {
		return
	} else if docErr != nil {
		if _, ok := docErr.(routing.ErrInvalidInputData); ok {
			http.Error(w, "invalid input data", 400)
		} else if _, ok := docErr.(routing.ErrNotAuthorize); ok {
			http.Error(w, "unauthorize", 401)
		} else if _, ok := docErr.(routing.ErrNotAuthenticate); ok {
			http.Error(w, "unauthorize", 401)
		} else {
			log.Println(docErr.Error())

			//print stack trace
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, true)
			log.Printf("%s\n", string(buf[0:stackSize]))

			http.Error(w, "internal error", 500)
		}
	}

	http.Error(w, "path not found", 404)
}

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

// Generate error page
func handleErrorCode(errorCode int, description string, w http.ResponseWriter) {
	http.Error(w, description, errorCode)

	// w.WriteHeader(errorCode)                    // set HTTP status code (example 404, 500)
	// w.Header().Set("Content-Type", "text/html") // clarify return type (MIME)
	// w.Write([]byte(fmt.Sprintf(
	// 	"<html><body><h1>Error %d</h1><p>%s</p></body></html>",
	// 	errorCode,
	// 	description)))
}
