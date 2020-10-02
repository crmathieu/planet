package storage
import (
	"github.com/crmathieu/planet/data"
	"net/http"
)

type DBStorage struct {
}

type DatabaseRepo struct {
	dsn string
	dbs DBStorage
}

var DBRepo DatabaseRepo

func (db DatabaseRepo)Init() error {
	return nil
}
func (db DatabaseRepo)UserGet(w http.ResponseWriter, r *http.Request, userid string)  {
}
func (db DatabaseRepo)UserAdd(w http.ResponseWriter, r *http.Request, user *data.USER)  {
}
func (db DatabaseRepo)UserDelete(w http.ResponseWriter, r *http.Request, userid string)  {
}
func (db DatabaseRepo)UserUpdate(w http.ResponseWriter, r *http.Request, userid string, user *data.USER)  {
}
func (db DatabaseRepo)GroupGet(w http.ResponseWriter, r *http.Request, groupname string)  {
}
func (db DatabaseRepo)GroupAdd(w http.ResponseWriter, r *http.Request, group *data.GROUP)  {
}
func (db DatabaseRepo)GroupDelete(w http.ResponseWriter, r *http.Request, groupname string)  {
}
func (db DatabaseRepo)GroupUpdate(w http.ResponseWriter, r *http.Request, groupname string, grpupd *data.GROUPUPD)  {
}

/*
func (db DBStorage)UserGet(w http.ResponseWriter, r *http.Request, userid string)  {
}
func (db DBStorage)UserAdd(w http.ResponseWriter, r *http.Request, user *data.USER)  {
}
func (db DBStorage)UserDelete(w http.ResponseWriter, r *http.Request, userid string)  {
}
func (db DBStorage)UserUpdate(w http.ResponseWriter, r *http.Request, userid string, user *data.USER)  {
}
func (db DBStorage)GroupGet(w http.ResponseWriter, r *http.Request, groupname string)  {
}
func (db DBStorage)GroupAdd(w http.ResponseWriter, r *http.Request, group *data.GROUP)  {
}
func (db DBStorage)GroupDelete(w http.ResponseWriter, r *http.Request, groupname string)  {
}
func (db DBStorage)GroupUpdate(w http.ResponseWriter, r *http.Request, groupname string, grpupd *data.GROUPUPD)  {
}
*/