package dbx

func (db *DBI) QueryStmt(tblName string, conds []Cond, options ...O) *queryStmt {
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}

	return &queryStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			conds: conds,
		},
		bys: opts.bys,
		limit: opts.limit,
	}
}

func (db *DBI) ListStmt(tblName string, conds []Cond, options ...O) *listStmt {
	return &listStmt{
		queryStmt: db.QueryStmt(tblName, conds, options...),
	}
}

func (db *DBI) SelectStmt(tblName string, fields []string, conds []Cond, options ...O) *selectStmt {
	return &selectStmt{
		listStmt: db.ListStmt(tblName, conds, options...),
		fields: fields,
	}
}

func (db *DBI) SqlStmt(tblName string, sql string, options ...O) *sqlStmt {
	return &sqlStmt{
		listStmt: db.ListStmt(tblName, nil, options...),
		sql: sql,
	}
}

func (db *DBI) InnerJoinStmt(tblName string, joinedTblName string, joinCond string, conds []Cond, options ...O) *innerJoinStmt {
	return &innerJoinStmt{
		listStmt: db.ListStmt(tblName, conds, options...),
		joinedTbl: joinedTblName,
		joinCond: joinCond,
	}
}

func (db *DBI) InsertStmt(tblName string) *insertStmt {
	return &insertStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
		},
	}
}

func (db *DBI) UpdateStmt(tblName string, conds []Cond, cols []string) *updateStmt {
	return &updateStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			conds: conds,
		},
		cols: cols,
	}
}

func (db *DBI) DeleteStmt(tblName string, conds []Cond) *deleteStmt {
	return &deleteStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			conds: conds,
		},
	}
}

func QueryStmt(tblName string, conds []Cond, options ...O) *queryStmt {
	db := getDefaultConnection()
	return db.QueryStmt(tblName, conds, options...)
}

func ListStmt(tblName string, conds []Cond, options ...O) *listStmt {
	db := getDefaultConnection()
	return db.ListStmt(tblName, conds, options...)
}

func SelectStmt(tblName string, fields []string, conds []Cond, options ...O) *selectStmt {
	db := getDefaultConnection()
	return db.SelectStmt(tblName, fields, conds, options...)
}

func SqlStmt(tblName string, sql string, options ...O) *sqlStmt {
	db := getDefaultConnection()
	return db.SqlStmt(tblName, sql, options...)
}

func InnerJoinStmt(tblName string, joinedTblName string, joinCond string, conds []Cond, options ...O) *innerJoinStmt {
	db := getDefaultConnection()
	return db.InnerJoinStmt(tblName, joinedTblName, joinCond, conds, options...)
}

func InsertStmt(tblName string) *insertStmt {
	db := getDefaultConnection()
	return db.InsertStmt(tblName)
}

func UpdateStmt(tblName string, conds []Cond, cols []string) *updateStmt {
	db := getDefaultConnection()
	return db.UpdateStmt(tblName, conds, cols)
}

func DeleteStmt(tblName string, conds []Cond) *deleteStmt {
	db := getDefaultConnection()
	return db.DeleteStmt(tblName, conds)
}

func VoidStmt() *voidStmt {
	return &voidStmt{}
}

// some re-usable handler
func (db *DBI) GetBy(tblName, colName string, colVal interface{}, res interface{}, options ...O) (bool, error) {
	stmt := db.QueryStmt(tblName, []Cond{Eq(colName, colVal)}, options...)
	r, err := stmt.Exec(res)
	if err != nil {
		return false, err
	}
	return r.(bool), nil
}

func (db *DBI) Get(tblName string, conds []Cond, res interface{}, options ...O) (bool, error) {
	stmt := db.QueryStmt(tblName, conds, options...)
	r, err := stmt.Exec(res)
	if err != nil {
		return false, err
	}
	return r.(bool), nil
}

func (db *DBI) GetOne(tblName string, conds []Cond, res interface{}, options ...O) (bool, error) {
	return db.Get(tblName, conds, res, options...)
}

func (db *DBI) List(tblName string, conds []Cond, res interface{}, options ...O) error {
	stmt := db.ListStmt(tblName, conds, options...)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Find(tblName string, conds []Cond, res interface{}, options ...O) error {
	return db.List(tblName, conds, res, options...)
}

func (db *DBI) InnerJoin(tblName string, joinedTblName string, joinCond string, conds []Cond, res interface{}, options ...O) error {
	stmt := db.InnerJoinStmt(tblName, joinedTblName, joinCond, conds, options...)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Select(tblName string, fields []string, conds []Cond, res interface{}, options ...O) error {
	stmt := db.SelectStmt(tblName, fields, conds, options...)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Insert(tblName string, vals interface{}) error {
	_, err := db.InsertStmt(tblName).Exec(vals)
	return err
}

func (db *DBI) Update(tblName string, conds []Cond, cols []string, vals interface{}) error {
	_, err := db.UpdateStmt(tblName, conds, cols).Exec(vals)
	return err
}

func (db *DBI) Delete(tblName string, conds []Cond, vals interface{}) error {
	_, err := db.DeleteStmt(tblName, conds).Exec(vals)
	return err
}

func (db *DBI) RunSQL(tblName string, sql string, res interface{}) error {
	stmt := db.SqlStmt(tblName, sql)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Iter(tblName string, conds []Cond, bean interface{}, options ...O) (<-chan interface{}) {
	stmt := db.QueryStmt(tblName, conds, options...)
	return stmt.Iter(bean)
}

func (db *DBI) Iterate(tblName string, conds []Cond, bean interface{}, it FnIterate, options ...O) error {
	stmt := db.QueryStmt(tblName, conds, options...)
	return stmt.Iterate(bean, it)
}

func GetBy(tblName, colName string, colVal interface{}, res interface{}, options ...O) (bool, error) {
	db := getDefaultConnection()
	return db.GetBy(tblName, colName, colVal, res, options...)
}

func Get(tblName string, conds []Cond, res interface{}, options ...O) (bool, error) {
	db := getDefaultConnection()
	return db.Get(tblName, conds, res, options...)
}

var GetOne = Get

func List(tblName string, conds []Cond, res interface{}, options ...O) error {
	db := getDefaultConnection()
	return db.List(tblName, conds, res, options...)
}

var Find = List

func InnerJoin(tblName string, joinedTblName string, joinCond string, conds []Cond, res interface{}, options ...O) error {
	db := getDefaultConnection()
	return db.InnerJoin(tblName, joinedTblName, joinCond, conds, res, options...)
}

func Select(tblName string, fields []string, conds []Cond, res interface{}, options ...O) error {
	db := getDefaultConnection()
	return db.Select(tblName, fields, conds, res, options...)
}

func Insert(tblName string, vals interface{}) error {
	db := getDefaultConnection()
	return db.Insert(tblName, vals)
}

func Update(tblName string, conds []Cond, cols []string, vals interface{}) error {
	db := getDefaultConnection()
	return db.Update(tblName, conds, cols, vals)
}

func Delete(tblName string, conds []Cond, vals interface{}) error {
	db := getDefaultConnection()
	return db.Delete(tblName, conds, vals)
}

func RunSQL(tblName string, sql string, res interface{}) error {
	db := getDefaultConnection()
	return db.RunSQL(tblName, sql, res)
}

func Iter(tblName string, conds []Cond, bean interface{}, options ...O) (<-chan interface{}) {
	db := getDefaultConnection()
	return db.Iter(tblName, conds, bean, options...)
}

func Iterate(tblName string, conds []Cond, bean interface{}, it FnIterate, options ...O) error {
	db := getDefaultConnection()
	return db.Iterate(tblName, conds, bean, it, options...)
}

// some statistic func
func (stmt *queryStmt) Count(bean interface{}, session ...*Session) (int64, error) {
	sess := stmt.createQuerySession(session)
	return sess.Count(bean)
}

func (stmt *queryStmt) Sum(bean interface{}, col string, session ...*Session) (float64, error) {
	sess := stmt.createQuerySession(session)
	return sess.Sum(bean, col)
}
