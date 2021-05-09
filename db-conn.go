// db connection related
package dbx

import (
	"github.com/rosbit/xorm"
)

type DBI struct {
	*xorm.Engine
}

var (
	DB *DBI // default db connection instance
)

// create a default instatnce of db connection for a driver
func CreateDBDriverConnection(driverName, dsn string, debug bool) (err error) {
	DB, err = CreateDriverDBInstance(driverName, dsn, debug)
	return
}

func getDefaultConnection() *DBI {
	if DB == nil {
		panic("please call CreateDBConnection(...) first")
	}
	return DB
}

// create an instance of connection for a driver
func CreateDriverDBInstance(driverName, dsn string, debug bool) (db *DBI, err error) {
	var dbInst *xorm.Engine
	dbInst, err = xorm.NewEngine(driverName, dsn)
	if err == nil {
		db = &DBI{dbInst}
		if debug {
			dbInst.ShowSQL(true)
		}
	}
	return
}
