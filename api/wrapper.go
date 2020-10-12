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
	if id != "" {
		(*Strg).UserGet(w, r, id)
	} else {
		data.ServerResponse(w, r, http.StatusUnprocessableEntity, data.STRING_DATA, "missing userid")
	}
}

// UserAddWrapper -------------------------------------------------------------
// POST method - extracts body and calls storage specific user add
// ----------------------------------------------------------------------------
func UserAddWrapper(w http.ResponseWriter, r *http.Request, dummy string) {
	// payload is in body
	defer r.Body.Close()

	plBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, fmt.Sprintf("Error extracting payload from request body: %s", err.Error()))
		return
	}

	var user data.USER
	err = json.Unmarshal(plBytes, &user)
	if err != nil {
		data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, fmt.Sprintf("Error unmarshalling payload - %s\nPAYLOAD=%s\n", err.Error(), string(plBytes)))
		return
	}
	(*Strg).UserAdd(w, r, &user)
}

// UserDeleteWrapper ----------------------------------------------------------
// DELETE method - calls storage specific user delete
// ----------------------------------------------------------------------------
func UserDeleteWrapper(w http.ResponseWriter, r *http.Request, id string) {
	if id != "" {
		(*Strg).UserDelete(w, r, id)
	} else {
		data.ServerResponse(w, r, http.StatusUnprocessableEntity, data.STRING_DATA, "missing userid")
	}
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
			data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, fmt.Sprintf("Error extracting payload from request body: %s", err.Error()))
			return
		}

		var user data.USER
		err = json.Unmarshal(plBytes, &user)
		if err != nil {
			data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, fmt.Sprintf("Error unmarshalling payload - %s\nPAYLOAD=%s\n", err.Error(), string(plBytes)))
			return
		}

		(*Strg).UserUpdate(w, r, id, &user)
	} else {
		data.ServerResponse(w, r, http.StatusUnprocessableEntity, data.STRING_DATA, "missing userid")
	}
}

//
// GROUP FAMILY
//

// GroupGetWrapper ----------------------------------------------------------
// GET method - calls storage specific group get
// ----------------------------------------------------------------------------
func GroupGetWrapper(w http.ResponseWriter, r *http.Request, id string) {
	if id != "" {
		(*Strg).GroupGet(w, r, id)
	} else {
		data.ServerResponse(w, r, http.StatusUnprocessableEntity, data.STRING_DATA, "missing groupname")
	}
}

// GroupDeleteWrapper ----------------------------------------------------------
// DELETE method - calls storage specific group delete
// ----------------------------------------------------------------------------
func GroupDeleteWrapper(w http.ResponseWriter, r *http.Request, id string) {
	if id != "" {
		(*Strg).GroupDelete(w, r, id)
	} else {
		data.ServerResponse(w, r, http.StatusUnprocessableEntity, data.STRING_DATA, "missing groupname")
	}
}

// GroupUpdateWrapper ----------------------------------------------------------
// PUT method - extracts body and calls storage specific group update
// ----------------------------------------------------------------------------
func GroupUpdateWrapper(w http.ResponseWriter, r *http.Request, id string) {
	// payload is in PUT body
	if id != "" {
		plBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, fmt.Sprintf("Error extracting payload from request body: %s", err.Error()))
			return
		}

		var grpupd data.GROUPUPD
		err = json.Unmarshal(plBytes, &grpupd.Members)
		if err != nil {
			data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, fmt.Sprintf("Error unmarshalling payload - %s\nPAYLOAD=%s\n", err.Error(), string(plBytes)))
			return
		}

		(*Strg).GroupUpdate(w, r, id, &grpupd)
	} else {
		data.ServerResponse(w, r, http.StatusUnprocessableEntity, data.STRING_DATA, "missing groupname")
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
		data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, fmt.Sprintf("Error extracting payload from request body: %s", err.Error()))
		return
	}

	var group data.GROUP
	err = json.Unmarshal(plBytes, &group)
	if err != nil {
		data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, fmt.Sprintf("Error unmarshalling payload - %s\nPAYLOAD=%s\n", err.Error(), string(plBytes)))
		return
	}
	(*Strg).GroupAdd(w, r, &group)
}
