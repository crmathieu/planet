package storage
import (
	"github.com/crmathieu/planet/data"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"sync"
)

type JsonStorage struct {
	Repo map[string]*data.USER `json:"repo"`
	Repogrp map[string][]string `json:"repogrp"`
} 

type JsonRepo struct {
	LastRepoName string
	uLock sync.RWMutex  // protects access to Repo maps
	js JsonStorage
}

type JsonLastRepo struct {
	LastRepo string `json:"lastrepo"`
}

var JRepo JsonRepo

// Init -----------------------------------------------------------------------
// json storage specific initialization
// ----------------------------------------------------------------------------
func (js JsonStorage) Init() error {

	if !loadRepoFromDisk() {
		JRepo.js = JsonStorage{}
	}
	return nil
}

// saveRepoToDisk -------------------------------------------------------------
func saveRepoToDisk() {
	var err error
	var body []byte
	if body, err = json.Marshal(JRepo.js); err == nil {
		// generate file name
		fn := fmt.Sprintf(time.Now().Format("20060102150405")) + ".json"
		fd, err1 := os.Create("storage/repository/"+fn)
		if err1 != nil {
			fmt.Println("** Could not save data to repository:",err1.Error(),"**")
			return
		}
		defer fd.Close()
		fh, err2 := os.OpenFile("storage/repository/last.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err2 != nil {
			fmt.Println("** Could not open index to repository:",err2.Error(),"**")
			return
		}
		defer fh.Close()
		var repoLast = JsonLastRepo{LastRepo: fn}
		out, _ := json.Marshal(repoLast) 
		_, err = fh.Write(out)
		if err != nil {
			fmt.Println("** Could not update index to repository:",err.Error(),"**")
			return
		}
		_, err = fd.Write(body)
		if err != nil {
			fmt.Println("** Could not save content to repository:",err.Error(),"**")
			return
		}
	} else {
		fmt.Println(err.Error())
	}
	// when everything goes well, remove previous storage:
	err = os.Remove("storage/repository/"+JRepo.LastRepoName) 
    if err != nil { 
		fmt.Println("** Could not remove older repository:",err.Error(),"**")
		return
    } 
}

// loadRepoFromDisk - load data -----------------------------------------------
func loadRepoFromDisk() bool {
    fnBytes, err := ioutil.ReadFile("storage/repository/last.txt")
	if err != nil {
		fmt.Println("** Could not update index to repository:",err.Error(),"**")
		return false
	}
	var fn JsonLastRepo
	err = json.Unmarshal(fnBytes, &fn)
	var body []byte
	fmt.Println("----",fn.LastRepo)
	body, err = ioutil.ReadFile("storage/repository/"+fn.LastRepo) //string(fn))
	if err !=  nil {
		fmt.Println("** Could not read last repository:",fn,"- error:",err.Error(),"**")
		return false
	}

	err = json.Unmarshal(body, &JRepo.js)
	if err != nil {
		fmt.Println("Error unmarshalling last repo backup:", fn, "- error", err.Error())
		return false
	} 
	JRepo.LastRepoName = string(fn.LastRepo)
	return true
}

// removes a userid from a group's member list --------------------------------
func removeItem(groupname, userid string) {
	for k, v := range JRepo.js.Repogrp[groupname] {
		if v == userid {
			JRepo.js.Repogrp[groupname] = append(JRepo.js.Repogrp[groupname][:k], JRepo.js.Repogrp[groupname][k+1:]...)
			break
		}
	}
}

// UserGet --------------------------------------------------------------------
// Get method - Returns the matching user record or 404
// ----------------------------------------------------------------------------
func (js JsonStorage) UserGet(w http.ResponseWriter, r *http.Request, userid string) {
	if userid != "" {
		JRepo.uLock.RLock()
		defer JRepo.uLock.RUnlock() 
		if u, ok := JRepo.js.Repo[userid]; ok {
			body, mok := json.Marshal(u)
			if mok == nil {
				data.ServerResponse(w, r, http.StatusOK, data.JSON_DATA, string(body))
			} else {
				data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, http.StatusText(http.StatusInternalServerError))
			}
		} else {
			data.ServerResponse(w, r, http.StatusNotFound, data.STRING_DATA, http.StatusText(http.StatusNotFound))
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, data.STRING_DATA, http.StatusText(http.StatusBadRequest))
	}
}


// UserAdd --------------------------------------------------------------------
// POST method - creates a new user record. Body contains USER information
// -----------------------------------------------------------------------	-----
func (js JsonStorage) UserAdd(w http.ResponseWriter, r *http.Request, user *data.USER) {

	JRepo.uLock.Lock()
	defer JRepo.uLock.Unlock() 
	if _, ok := JRepo.js.Repo[user.UID]; !ok {
		// userid not present in Repo
		JRepo.js.Repo[user.UID] = user
		for _, v := range user.Groups {
			JRepo.js.Repogrp[v] = append(JRepo.js.Repogrp[v], user.UID) 
		}
		data.ServerResponse(w, r, http.StatusOK, data.STRING_DATA, "User added successfully")
	} else {
		data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, "Error user "+user.UID+" already exists")
	}
}

// UserDelete -----------------------------------------------------------------
// DELETE method - deletes an existing user record
// ----------------------------------------------------------------------------
func (js JsonStorage) UserDelete(w http.ResponseWriter, r *http.Request, userid string)  {
	if userid != "" {
		JRepo.uLock.Lock()
		defer JRepo.uLock.Unlock() 
		if u, ok := JRepo.js.Repo[userid]; ok {
			for _, v := range u.Groups {
				removeItem(v, userid)
			}
			delete(JRepo.js.Repo, userid)
			data.ServerResponse(w, r, http.StatusOK, data.STRING_DATA, "User "+userid+" was sucessfully deleted")
		} else {
			data.ServerResponse(w, r, http.StatusNotFound, data.STRING_DATA, http.StatusText(http.StatusNotFound))
		}
		
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, data.STRING_DATA, http.StatusText(http.StatusBadRequest))
	}
}

// UserUpdate -----------------------------------------------------------------
// PUT method - updates an existing user record
// ----------------------------------------------------------------------------
func (js JsonStorage) UserUpdate(w http.ResponseWriter, r *http.Request, userid string, user *data.USER)  {

	JRepo.uLock.Lock()
	defer JRepo.uLock.Unlock() 

	if u, ok := JRepo.js.Repo[userid]; ok {
		// userid not present in Repo
		for _, v := range u.Groups {
			removeItem(v, userid)
		}
		JRepo.js.Repo[userid] = user
		for _, v := range user.Groups {
			JRepo.js.Repogrp[v] = append(JRepo.js.Repogrp[v], user.UID) 
		}
		data.ServerResponse(w, r, http.StatusOK, data.STRING_DATA, "User "+userid+" was updated successfully")
	} else {
		data.ServerResponse(w, r, http.StatusNotFound, data.STRING_DATA, http.StatusText(http.StatusNotFound))
	}
}

// GroupGet -------------------------------------------------------------------
// GET method - returns the group's members list
// ----------------------------------------------------------------------------
func (js JsonStorage) GroupGet(w http.ResponseWriter, r *http.Request, groupname string)  {
	JRepo.uLock.RLock()
	defer JRepo.uLock.RUnlock() 
	if grp, ok := JRepo.js.Repogrp[groupname]; ok {
		body, mok := json.Marshal(grp)
		if mok == nil {
			data.ServerResponse(w, r, http.StatusOK, data.JSON_DATA, string(body))
		} else {
			data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, http.StatusText(http.StatusInternalServerError))
		}
	} else {
		data.ServerResponse(w, r, http.StatusNotFound, data.STRING_DATA, http.StatusText(http.StatusNotFound))
	}
}

// GroupAdd -------------------------------------------------------------------
// POST method - Creates an empty group
// ----------------------------------------------------------------------------
func (js JsonStorage) GroupAdd(w http.ResponseWriter, r *http.Request, group *data.GROUP)  {
	if group.Gname != "" {
		JRepo.uLock.Lock()
		defer JRepo.uLock.Unlock() 
		if _, ok := JRepo.js.Repogrp[group.Gname]; !ok {
			JRepo.js.Repogrp[group.Gname] = []string{}
			data.ServerResponse(w, r, http.StatusOK, data.STRING_DATA, "Group "+group.Gname+" added successfully")
		} else {
			data.ServerResponse(w, r, http.StatusInternalServerError, data.STRING_DATA, "Group "+group.Gname+" already exists")
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, data.STRING_DATA, http.StatusText(http.StatusBadRequest))
	}
}

// GroupDelete ----------------------------------------------------------------
// DELETED method - removes group and references to it
// ----------------------------------------------------------------------------
func (js JsonStorage) GroupDelete(w http.ResponseWriter, r *http.Request, groupname string)  {
	if groupname != "" {
		JRepo.uLock.Lock()
		defer JRepo.uLock.Unlock() 
		if _, ok := JRepo.js.Repogrp[groupname]; ok {
			// if groupname exists, remove all references to it for each member in its list
			if v, ok2 := JRepo.js.Repogrp[groupname]; ok2 {
				removeAllReference(&v, groupname)
			}
			// then remove group name from map
			delete(JRepo.js.Repogrp, groupname)
			data.ServerResponse(w, r, http.StatusOK, data.STRING_DATA, "Group "+groupname+" successfully deleted")
		} else {
			data.ServerResponse(w, r, http.StatusNotFound, data.STRING_DATA, http.StatusText(http.StatusNotFound))
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, data.STRING_DATA, http.StatusText(http.StatusBadRequest))
	}
}

// removes all group name reference for each user in group list ---------------
func removeAllReference(list *[]string, groupname string) {
	for _, v := range *list {
		removeReference(v, groupname)
	}
}


// removes a group name from a user record ------------------------------------
func removeReference(userid string, groupname string) {
	user := JRepo.js.Repo[userid]
	for k, gn := range(user.Groups) {
		if gn == groupname {
			JRepo.js.Repo[user.UID].Groups = append(JRepo.js.Repo[user.UID].Groups[:k], JRepo.js.Repo[user.UID].Groups[k+1:]...)
		}
	}
}

// adds a group name to a user record -----------------------------------------
func addReference(userid, groupname string) {
	if _, ok := JRepo.js.Repo[userid]; ok {
		JRepo.js.Repo[userid].Groups = append(JRepo.js.Repo[userid].Groups, groupname)
	}
}

// GroupUpdate ----------------------------------------------------------------
// PUT method - updates a group members list. New list is in body as array of
//              strings
// ----------------------------------------------------------------------------
func (js JsonStorage)GroupUpdate(w http.ResponseWriter, r *http.Request, groupname string, grpupd *data.GROUPUPD)  {
	if groupname != "" {
		JRepo.uLock.Lock()
		defer JRepo.uLock.Unlock() 
		if oldmembers, ok := JRepo.js.Repogrp[groupname]; ok {
			var aux = make(map[string]int8)
			// create a map of new set of members
			for _, uid := range grpupd.Members {
				aux[uid] = 0
			}
			// circle through old set and remove group name from each member not in new list
			for _, userid := range (oldmembers) {
				if _, ok := aux[userid]; !ok {
					removeReference(userid, groupname)
				} else {
					aux[userid] = 1
				}
			}
			// finally add group name for users not already in oldmember list
			for k, v := range aux {
				if v == 0 {
					addReference(k, groupname)
				}
			}
			// ... and replace old set with new one
			JRepo.js.Repogrp[groupname] = grpupd.Members
			data.ServerResponse(w, r, http.StatusOK, data.STRING_DATA, "Group "+groupname+" was successfully updated")
		} else {
			data.ServerResponse(w, r, http.StatusNotAcceptable, data.STRING_DATA, http.StatusText(http.StatusNotAcceptable))
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, data.STRING_DATA, http.StatusText(http.StatusBadRequest))
	}
}
