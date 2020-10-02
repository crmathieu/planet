package data
import (
    "net/http"
)
type USER struct {
	Fname string `json:"first_name"`     
    Lname string `json:"last_name"`     
    UID    string `json:"userid"`     
    Groups []string `json:"groups"` 
}

type GROUP struct {
	Gname string `json:"group_name"`     
}

type GROUPUPD struct {
    Members []string `json:""`
}

func ServerResponse (w http.ResponseWriter, r *http.Request, serverCode int, body []byte) {
    w.WriteHeader(serverCode)
    w.Write(body)
}
