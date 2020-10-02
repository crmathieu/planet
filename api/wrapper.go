package api
import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/crmathieu/planet/data"
	"github.com/crmathieu/planet/storage"
	"io/ioutil"
)

var Strg *storage.UserStorage

// InitStorage ----------------------------------------------------------------
func InitStorage(storageType string) error {
	var err error
	Strg, err = storage.OpenStore(storageType)
	return err
}

// CloseStorage ---------------------------------------------------------------
func CloseStorage(storageType string) {
	storage.CloseStore(storageType)
}

//
// USER FAMILY
//

// UserGetWrapper -------------------------------------------------------------
// Get method - calls storage specific user get
// ----------------------------------------------------------------------------
func UserGetWrapper(w http.ResponseWriter, r *http.Request, id string) {
	(*Strg).UserGet(w, r, id)
}

// UserAddWrapper -------------------------------------------------------------
// POST method - extracts body and calls storage specific user add
// ----------------------------------------------------------------------------
func UserAddWrapper(w http.ResponseWriter, r *http.Request, dummy string) {
	// payload is in body
	defer r.Body.Close()

	plBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data.ServerResponse(w, r, http.StatusInternalServerError, []byte(fmt.Sprintf("Error extracting payload from request body: %s", err.Error())))
		return
	}

	var user data.USER
	err = json.Unmarshal(plBytes, &user)
	if err != nil {
		data.ServerResponse(w, r, http.StatusInternalServerError, []byte(fmt.Sprintf("Error unmarshalling payload - %s\nPAYLOAD=%s\n", err.Error(), string(plBytes))))
		return
	}
	(*Strg).UserAdd(w, r, &user)
}

// UserDeleteWrapper ----------------------------------------------------------
// DELETE method - calls storage specific user delete
// ----------------------------------------------------------------------------
func UserDeleteWrapper(w http.ResponseWriter, r *http.Request, id string) {
	(*Strg).UserDelete(w, r, id)
}

// UserUpdateWrapper ----------------------------------------------------------
// PUT method - extracts body and calls storage specific user update
// ----------------------------------------------------------------------------
func UserUpdateWrapper(w http.ResponseWriter, r *http.Request, id string) {
	// payload is in body
	defer r.Body.Close()

	if id != "" {
		plBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			data.ServerResponse(w, r, http.StatusInternalServerError, []byte(fmt.Sprintf("Error extracting payload from request body: %s", err.Error())))
			return
		}

		var user data.USER
		err = json.Unmarshal(plBytes, &user)
		if err != nil {
			data.ServerResponse(w, r, http.StatusInternalServerError, []byte(fmt.Sprintf("Error unmarshalling payload - %s\nPAYLOAD=%s\n", err.Error(), string(plBytes))))
			return
		}

		(*Strg).UserUpdate(w, r, id, &user)
	} else {
		data.ServerResponse(w, r, http.StatusUnprocessableEntity, []byte("missing userid not found"))
	}
}

//
// GROUP FAMILY
//

// GroupGetWrapper ----------------------------------------------------------
// GET method - calls storage specific group get
// ----------------------------------------------------------------------------
func GroupGetWrapper(w http.ResponseWriter, r *http.Request, id string) {
	(*Strg).GroupGet(w, r, id)
}

// GroupDeleteWrapper ----------------------------------------------------------
// DELETE method - calls storage specific group delete
// ----------------------------------------------------------------------------
func GroupDeleteWrapper(w http.ResponseWriter, r *http.Request, id string) {
	(*Strg).GroupDelete(w, r, id)
}

// GroupUpdateWrapper ----------------------------------------------------------
// PUT method - extracts body and calls storage specific group update
// ----------------------------------------------------------------------------
func GroupUpdateWrapper(w http.ResponseWriter, r *http.Request, id string) {
	// payload is in PUT body
	if id != "" {
		plBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			data.ServerResponse(w, r, http.StatusInternalServerError, []byte(fmt.Sprintf("Error extracting payload from request body: %s", err.Error())))
			return
		}

		var grpupd data.GROUPUPD
		err = json.Unmarshal(plBytes, &grpupd.Members)
		if err != nil {
			data.ServerResponse(w, r, http.StatusInternalServerError, []byte(fmt.Sprintf("Error unmarshalling payload - %s\nPAYLOAD=%s\n", err.Error(), string(plBytes))))
			return
		}

		(*Strg).GroupUpdate(w, r, id, &grpupd)
	} else {
		data.ServerResponse(w, r, http.StatusUnprocessableEntity, []byte("missing groupname"))
	}
}

// GroupAddWrapper ----------------------------------------------------------
// POST method - extracts body and calls storage specific group add
// ----------------------------------------------------------------------------
func GroupAddWrapper(w http.ResponseWriter, r *http.Request, dummy string) {
	// payload is in POST body
	defer r.Body.Close()

	plBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data.ServerResponse(w, r, http.StatusInternalServerError, []byte(fmt.Sprintf("Error extracting payload from request body: %s", err.Error())))
		return
	}

	var group data.GROUP
	err = json.Unmarshal(plBytes, &group)
	if err != nil {
		data.ServerResponse(w, r, http.StatusInternalServerError, []byte(fmt.Sprintf("Error unmarshalling payload - %s\nPAYLOAD=%s\n", err.Error(), string(plBytes))))
		return
	}
	(*Strg).GroupAdd(w, r, &group)
}
