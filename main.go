package main

import (
	"net/http"
	"io"
	"github.com/crmathieu/planet/controller"
	"github.com/crmathieu/planet/api"
	"strings"
)

/*const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)*/

var routes = map[string]func(w http.ResponseWriter, r *http.Request, id string) {
	"users":    controller.UsersEndpoints,
	"groups":   controller.GroupsEndpoints,
}

// router ---------------------------------------------------------------------
// determines endpoint family
// ----------------------------------------------------------------------------
func router(w http.ResponseWriter, r *http.Request) {
	tokens := strings.Split(r.URL.String(), "/") 
	id := ""
	if len(tokens) > 2 {
		// 3rd token id the id
		id = tokens[2]
	} 
	// 2nd token is endp
	if fc, ok := routes[tokens[1]]; ok {
		// endpoint family recognized. Call it!
		fc(w, r, id)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, r.URL.Path + ": " +http.StatusText(http.StatusBadRequest))
	}
}

func init() {
	api.InitStorage("JSON")
}

func main() {

	http.HandleFunc("/", router)

	// CatchSignals will catch SIGINT and SIGTERM signals to determine shutdown time
	CatchSignals(&http.Server{Addr: ":80", Handler: nil})

	// storage specific close
	api.CloseStorage("JSON")
}

