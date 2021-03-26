package dbx

func (db *DBI) NewQueryStmt(tblName string, conds []Cond, options ...O) *queryStmt {
	var bys []By
	if len(options) > 0 {
		bys = options[0].Bys
	}
	return &queryStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			conds: conds,
		},
		bys: bys,
	}
}

func (db *DBI) NewListStmt(tblName string, conds []Cond, options ...O) *listStmt {
	var limit Limit
	if len(options) > 0 {
		limit = options[0].Count
	}
	return &listStmt{
		queryStmt: db.NewQueryStmt(tblName, conds, options...),
		limit: limit,
	}
}

func (db *DBI) NewSelectStmt(tblName string, fields []string, conds []Cond, options ...O) *selectStmt {
	return &selectStmt{
		listStmt: db.NewListStmt(tblName, conds, options...),
		fields: fields,
	}
}

func (db *DBI) NewSqlStmt(tblName string, sql string, options ...O) *sqlStmt {
	return &sqlStmt{
		listStmt: db.NewListStmt(tblName, nil, options...),
		sql: sql,
	}
}

func (db *DBI) NewInnerJoinStmt(tblName string, joinedTblName string, joinCond string, conds []Cond, options ...O) *innerJoinStmt {
	return &innerJoinStmt{
		listStmt: db.NewListStmt(tblName, conds, options...),
		joinedTbl: joinedTblName,
		joinCond: joinCond,
	}
}

func (db *DBI) NewInsertStmt(tblName string) *insertStmt {
	return &insertStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
		},
	}
}

func (db *DBI) NewUpdateStmt(tblName string, conds []Cond, cols []string) *updateStmt {
	return &updateStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			conds: conds,
		},
		cols: cols,
	}
}

func (db *DBI) NewDeleteStmt(tblName string, conds []Cond) *deleteStmt {
	return &deleteStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			conds: conds,
		},
	}
}

func NewQueryStmt(tblName string, conds []Cond, options ...O) *queryStmt {
	db := getDefaultConnection()
	return db.NewQueryStmt(tblName, conds, options...)
}

func NewListStmt(tblName string, conds []Cond, options ...O) *listStmt {
	db := getDefaultConnection()
	return db.NewListStmt(tblName, conds, options...)
}

func NewSelectStmt(tblName string, fields []string, conds []Cond, options ...O) *selectStmt {
	db := getDefaultConnection()
	return db.NewSelectStmt(tblName, fields, conds, options...)
}

func NewSqlStmt(tblName string, sql string, options ...O) *sqlStmt {
	db := getDefaultConnection()
	return db.NewSqlStmt(tblName, sql, options...)
}

func NewInnerJoinStmt(tblName string, joinedTblName string, joinCond string, conds []Cond, options ...O) *innerJoinStmt {
	db := getDefaultConnection()
	return db.NewInnerJoinStmt(tblName, joinedTblName, joinCond, conds, options...)
}

func NewInsertStmt(tblName string) *insertStmt {
	db := getDefaultConnection()
	return db.NewInsertStmt(tblName)
}

func NewUpdateStmt(tblName string, conds []Cond, cols []string) *updateStmt {
	db := getDefaultConnection()
	return db.NewUpdateStmt(tblName, conds, cols)
}

func NewDeleteStmt(tblName string, conds []Cond) *deleteStmt {
	db := getDefaultConnection()
	return db.NewDeleteStmt(tblName, conds)
}

func NewVoidStmt() *voidStmt {
	return &voidStmt{}
}

// some re-usable handler
func (db *DBI) GetBy(tblName, colName string, colVal interface{}, res interface{}) (bool, error) {
	stmt := db.NewQueryStmt(tblName, []Cond{Eq(colName, colVal)})
	res, err := stmt.Exec(res)
	if err != nil {
		return false, err
	}
	return res.(bool), err
}

func (db *DBI) Get(tblName string, conds []Cond, res interface{}) (bool, error) {
	stmt := db.NewQueryStmt(tblName, conds)
	res, err := stmt.Exec(res)
	if err != nil {
		return false, err
	}
	return res.(bool), err
}

func (db *DBI) Find(tblName string, conds []Cond, res interface{}, count ...Limit) error {
	var options O
	if len(count) > 0 {
		options.Count = count[0]
	}
	stmt := db.NewListStmt(tblName, conds, options)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Select(tblName string, fields []string, conds []Cond, res interface{}) error {
	stmt := db.NewSelectStmt(tblName, fields, conds)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Insert(tblName string, vals interface{}) error {
	_, err := db.NewInsertStmt(tblName).Exec(vals)
	return err
}

func (db *DBI) Update(tblName string, conds []Cond, cols []string, vals interface{}) error {
	_, err := db.NewUpdateStmt(tblName, conds, cols).Exec(vals)
	return err
}

func (db *DBI) Delete(tblName string, conds []Cond, vals interface{}) error {
	_, err := db.NewDeleteStmt(tblName, conds).Exec(vals)
	return err
}

func (db *DBI) RunSQL(tblName string, sql string, res interface{}) error {
	stmt := db.NewSqlStmt(tblName, sql)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Iter(tblName string, conds []Cond, bean interface{}) (<-chan interface{}) {
	stmt := db.NewQueryStmt(tblName, conds)
	return stmt.Iter(bean)
}

func (db *DBI) Iterate(tblName string, conds []Cond, bean interface{}, it FnIterate) error {
	stmt := db.NewQueryStmt(tblName, conds)
	return stmt.Iterate(bean, it)
}

func GetBy(tblName, colName string, colVal interface{}, res interface{}) (bool, error) {
	db := getDefaultConnection()
	return db.GetBy(tblName, colName, colVal, res)
}

func Get(tblName string, conds []Cond, res interface{}) (bool, error) {
	db := getDefaultConnection()
	return db.Get(tblName, conds, res)
}

func Find(tblName string, conds []Cond, res interface{}, count ...Limit) error {
	db := getDefaultConnection()
	return db.Find(tblName, conds, res, count...)
}

func Select(tblName string, fields []string, conds []Cond, res interface{}) error {
	db := getDefaultConnection()
	return db.Select(tblName, fields, conds, res)
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

func Iter(tblName string, conds []Cond, bean interface{}) (<-chan interface{}) {
	db := getDefaultConnection()
	return db.Iter(tblName, conds, bean)
}

func Iterate(tblName string, conds []Cond, bean interface{}, it FnIterate) error {
	db := getDefaultConnection()
	return db.Iterate(tblName, conds, bean, it)
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
