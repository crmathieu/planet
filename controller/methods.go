package controller
import (
	"net/http"
	"github.com/crmathieu/planet/api"
)

const (
	USERS_FAMILY = 0
	GROUPS_FAMILY = 1
)

var HttpMethod = map[string][]func(http.ResponseWriter, *http.Request, string){
	"GET": 		{api.UserGetWrapper, api.GroupGetWrapper},
	"DELETE":	{api.UserDeleteWrapper, api.GroupDeleteWrapper},
	"POST":		{api.UserAddWrapper, api.GroupAddWrapper},
	"PUT":		{api.UserUpdateWrapper, api.GroupUpdateWrapper},
} 

func callEndPoint(family int, w http.ResponseWriter, r *http.Request, id string) {
	if entry, ok := HttpMethod[r.Method]; ok {
		entry[family](w, r, id)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method " +r.Method+ " not allowed for users endpoints"))
	}
}

func UsersEndpoints(w http.ResponseWriter, r *http.Request, id string) { 
	callEndPoint(USERS_FAMILY, w, r, id)
}

func GroupsEndpoints(w http.ResponseWriter, r *http.Request, id string) { 
	callEndPoint(GROUPS_FAMILY, w, r, id)
}

