package storage
import (
	"github.com/crmathieu/planet/data"
	"net/http"
	"errors"
)

type UserStorage interface {
	Init() error
	UserGet(w http.ResponseWriter, r *http.Request, userid string) 
	UserAdd(w http.ResponseWriter, r *http.Request, user *data.USER) 
	UserDelete(w http.ResponseWriter, r *http.Request, userid string) 
	UserUpdate(w http.ResponseWriter, r *http.Request, userid string, user *data.USER) 

	GroupGet(w http.ResponseWriter, r *http.Request, groupname string) 
	GroupAdd(w http.ResponseWriter, r *http.Request, group *data.GROUP) 
	GroupDelete(w http.ResponseWriter, r *http.Request, groupname string) 
	GroupUpdate(w http.ResponseWriter, r *http.Request, groupname string, grpupd *data.GROUPUPD) 
}


// OpenStore ------------------------------------------------------------------
// open store based on requested type of storage 
// ----------------------------------------------------------------------------
func OpenStore(storageType string) (*UserStorage, error) {
	var store UserStorage
	switch storageType {
		case "JSON": store = JsonRepo{};break
		case "DB": store = DatabaseRepo{};break
		default: return nil, errors.New("Unknown storage type")
	}
	store.Init()
	return &store, nil
}

// CloseStore -----------------------------------------------------------------
// close store in use
// ----------------------------------------------------------------------------
func CloseStore(storageType string) {
	switch storageType {
		case "JSON": saveRepoToDisk();break
	}
}