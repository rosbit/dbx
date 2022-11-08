package dbx

func SyncTable(pTblStruct interface{}, tblName ...string) error {
	db := getDefaultConnection()
	return db.SyncTable(pTblStruct, tblName...)
}

func (db *DBI) SyncTable(pTblStruct interface{}, tblName ...string) error {
	if len(tblName) > 0 && len(tblName[0]) > 0 {
		return db.Engine.Table(tblName[0]).Sync2(pTblStruct)
	}
	return db.Engine.Sync2(pTblStruct)
}
