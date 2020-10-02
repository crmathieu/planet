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

var JRepo JsonRepo

func (js JsonRepo) Init() error {

//func (jr JsonRepo) Init() error {
/*	JsonRepo.Repo = map[string]*data.USER{"johnf": {"john","flemming", "johnf", []string{"users"}},
										 "joes": {"joseph","smith", "joes", []string{"users","admin"}},}

	JsonRepo.Repogrp = map[string][]string{
		"users":[]string{"johnf","joes"},
		"admin":[]string{"joes"},
	}
*/
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
		_, err = fh.Write([]byte(fn))
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
		fmt.Println("** Could remove older repository:",err.Error(),"**")
		return
    } 
}

// loadRepoFromDisk - load data -----------------------------------------------
func loadRepoFromDisk() bool {
    fn, err := ioutil.ReadFile("storage/repository/last.txt")
	if err != nil {
		fmt.Println("** Could not update index to repository:",err.Error(),"**")
		return false
	}
	var body []byte
	body, err = ioutil.ReadFile("storage/repository/"+string(fn))
	if err !=  nil {
		fmt.Println("** Could not read last repository:",fn,"- error:",err.Error(),"**")
		return false
	}

	err = json.Unmarshal(body, &JRepo.js)
	if err != nil {
		fmt.Println("Error unmarshalling last repo backup:", fn, "- error", err.Error())
		return false
	} 
	JRepo.LastRepoName = string(fn)
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
func (jr JsonRepo) UserGet(w http.ResponseWriter, r *http.Request, userid string) {
	if userid != "" {
		JRepo.uLock.RLock()
		defer JRepo.uLock.RUnlock() 
		if u, ok := JRepo.js.Repo[userid]; ok {
			body, mok := json.Marshal(u)
			if mok == nil {
				data.ServerResponse(w, r, http.StatusOK, body)
			} else {
				data.ServerResponse(w, r, http.StatusInternalServerError, []byte(http.StatusText(http.StatusInternalServerError)))
			}
		} else {
			data.ServerResponse(w, r, http.StatusNotFound, []byte(http.StatusText(http.StatusNotFound)))
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, []byte(http.StatusText(http.StatusBadRequest)))
	}
}


// UserAdd --------------------------------------------------------------------
// POST method - creates a new user record. Body contains USER information
// -----------------------------------------------------------------------	-----
func (jr JsonRepo) UserAdd(w http.ResponseWriter, r *http.Request, user *data.USER) {

	JRepo.uLock.Lock()
	defer JRepo.uLock.Unlock() 
	if _, ok := JRepo.js.Repo[user.UID]; !ok {
		// userid not present in Repo
		JRepo.js.Repo[user.UID] = user
		for _, v := range user.Groups {
			JRepo.js.Repogrp[v] = append(JRepo.js.Repogrp[v], user.UID) 
		}
		data.ServerResponse(w, r, http.StatusOK, []byte("User added successfully"))
	} else {
		data.ServerResponse(w, r, http.StatusInternalServerError, []byte("Error user "+user.UID+" already exists"))
	}
}

// UserDelete -----------------------------------------------------------------
// DELETE method - deletes an existing user record
// ----------------------------------------------------------------------------
func (jr JsonRepo) UserDelete(w http.ResponseWriter, r *http.Request, userid string)  {
	if userid != "" {
		JRepo.uLock.Lock()
		defer JRepo.uLock.Unlock() 
		if u, ok := JRepo.js.Repo[userid]; ok {
			for _, v := range u.Groups {
				removeItem(v, userid)
			}
			delete(JRepo.js.Repo, userid)
			data.ServerResponse(w, r, http.StatusOK, []byte("User "+userid+" was sucessfully deleted"))
		} else {
			data.ServerResponse(w, r, http.StatusNotFound, []byte(http.StatusText(http.StatusNotFound)))
		}
		
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, []byte(http.StatusText(http.StatusBadRequest)))
	}
}

// UserUpdate -----------------------------------------------------------------
// PUT method - updates an existing user record
// ----------------------------------------------------------------------------
func (jr JsonRepo) UserUpdate(w http.ResponseWriter, r *http.Request, userid string, user *data.USER)  {

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
		data.ServerResponse(w, r, http.StatusOK, []byte("User "+userid+" was updated successfully"))
	} else {
		data.ServerResponse(w, r, http.StatusNotFound, []byte(http.StatusText(http.StatusNotFound)))
	}
}

// GroupGet -------------------------------------------------------------------
// GET method - returns the group's members list
// ----------------------------------------------------------------------------
func (jr JsonRepo) GroupGet(w http.ResponseWriter, r *http.Request, groupname string)  {
	if groupname != "" {
		JRepo.uLock.RLock()
		defer JRepo.uLock.RUnlock() 
		if grp, ok := JRepo.js.Repogrp[groupname]; ok {
			body, mok := json.Marshal(grp)
			if mok == nil {
				data.ServerResponse(w, r, http.StatusOK, body)
			} else {
				data.ServerResponse(w, r, http.StatusInternalServerError, []byte(http.StatusText(http.StatusInternalServerError)))
			}
		} else {
			data.ServerResponse(w, r, http.StatusNotFound, []byte(http.StatusText(http.StatusNotFound)))
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, []byte(http.StatusText(http.StatusBadRequest)))
	}
}

// GroupAdd -------------------------------------------------------------------
// POST method - Creates an empty group
// ----------------------------------------------------------------------------
func (jr JsonRepo) GroupAdd(w http.ResponseWriter, r *http.Request, group *data.GROUP)  {
	if group.Gname != "" {
		JRepo.uLock.Lock()
		defer JRepo.uLock.Unlock() 
		if _, ok := JRepo.js.Repogrp[group.Gname]; !ok {
			JRepo.js.Repogrp[group.Gname] = []string{}
			data.ServerResponse(w, r, http.StatusOK, []byte("Group "+group.Gname+" added successfully"))
		} else {
			data.ServerResponse(w, r, http.StatusInternalServerError, []byte("Group "+group.Gname+" already exists"))
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, []byte(http.StatusText(http.StatusBadRequest)))
	}
}

// GroupDelete ----------------------------------------------------------------
// DELETED method - removes group and references to it
// ----------------------------------------------------------------------------
func (jr JsonRepo) GroupDelete(w http.ResponseWriter, r *http.Request, groupname string)  {
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
			data.ServerResponse(w, r, http.StatusOK, []byte("Group "+groupname+" successfully deleted"))
		} else {
			data.ServerResponse(w, r, http.StatusNotFound, []byte(http.StatusText(http.StatusNotFound)))
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, []byte(http.StatusText(http.StatusBadRequest)))
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
func (jr JsonRepo)GroupUpdate(w http.ResponseWriter, r *http.Request, groupname string, grpupd *data.GROUPUPD)  {
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
			data.ServerResponse(w, r, http.StatusOK, []byte("Group "+groupname+" was successfully updated"))
		} else {
			data.ServerResponse(w, r, http.StatusNotAcceptable, []byte(http.StatusText(http.StatusNotAcceptable)))
		}
	} else {
		data.ServerResponse(w, r, http.StatusBadRequest, []byte(http.StatusText(http.StatusBadRequest)))
	}
}
