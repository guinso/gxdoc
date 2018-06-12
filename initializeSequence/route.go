package initializeSequence

import (
	"log"
	"net/http"
	"strings"
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
	log.Println("Serving URL: " + r.URL.Path)

	//proceed to development if path matched
	if enableDevMode == true && (strings.HasPrefix(urlPath, devStartURL+"/") ||
		strings.Compare(urlPath, devStartURL) == 0) {

		urlPath = urlPath[len(devStartURL):]
		if strings.HasPrefix(urlPath, "/") {
			urlPath = urlPath[1:]
		}

		//if start with "api/" direct to REST handler
		if strings.HasPrefix(urlPath, "api/") || strings.Compare(urlPath, "api") == 0 {
			log.Println("route handle by development API")
			if strings.HasPrefix(urlPath, "api/") {
				routeDevelopmentAPI(w, r, urlPath[4:]) //trim: 'api/'
			} else {
				routeDevelopmentAPI(w, r, "")
			}
		} else {
			log.Println("route handle by development static file")
			//other wise, lets read a file path and display to client
			http.ServeFile(w, r, "./"+devStaticPath+"/"+urlPath)
		}

	} else {
		//if start with "api/" direct to REST handler
		if strings.HasPrefix(urlPath, "api/") || strings.Compare(urlPath, "api") == 0 {
			log.Println("route handle by production API")
			if strings.HasPrefix(urlPath, "api/") {
				routeProductionAPI(w, r, urlPath[4:]) //trim: "api/"
			} else {
				routeProductionAPI(w, r, "")
			}

		} else {
			log.Println("route handle by production stati file")
			//other wise, lets read a file path and display to client
			http.ServeFile(w, r, "./"+staticPath+"/"+urlPath)
		}
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
