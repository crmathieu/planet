package controller
import (
	"net/http"
	"github.com/crmathieu/planet/api"
)

var UsersHttpMethod = map[string]func(http.ResponseWriter, *http.Request, string){
	"GET": 		api.UserGetWrapper,
	"DELETE":	api.UserDeleteWrapper,
	"POST":		api.UserAddWrapper,
	"PUT":		api.UserUpdateWrapper,
} 

var GroupsHttpMethod = map[string]func(http.ResponseWriter, *http.Request, string){
	"GET": 		api.GroupGetWrapper,
	"DELETE":	api.GroupDeleteWrapper,
	"POST":		api.GroupAddWrapper,
	"PUT":		api.GroupUpdateWrapper,
} 

func UsersEndpoints(w http.ResponseWriter, r *http.Request, id string) { 

	if verb, ok := UsersHttpMethod[r.Method]; ok {
		verb(w, r, id)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method " +r.Method+ " not allowed for users endpoints"))
	}
}

func GroupsEndpoints(w http.ResponseWriter, r *http.Request, id string) { 

	if verb, ok := GroupsHttpMethod[r.Method]; ok {
		verb(w, r, id)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method " +r.Method+ " not allowed for groups endpoints"))
	}
}

