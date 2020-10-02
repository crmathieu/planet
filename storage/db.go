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

func (db DBStorage)Init() error {
	return nil
}
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
