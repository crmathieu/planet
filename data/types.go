package data
import (
    "net/http"
    "fmt"
)

const (
    JSON_DATA = 0
    STRING_DATA = 1
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

func ServerResponse (w http.ResponseWriter, r *http.Request, serverCode int, responseType int, body string) {
    var quote = ``
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(serverCode)
    if responseType == STRING_DATA {
        quote = `"`
    }
    w.Write([]byte(fmt.Sprintf(`{"status": %d, "response": %v%v%v}`, serverCode, quote, body, quote)))
}
