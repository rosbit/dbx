package dbx

func (db *DBI) QueryStmt(tblName string, conds []Cond, options ...O) *queryStmt {
	opts := getOptions(options...)

	return &queryStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			conds: conds,
			session: opts.session,
		},
		bys: opts.bys,
		limit: opts.limit,
		selection: opts.selection,
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

func (db *DBI) joinStmt(tblName string, joinedTblName string, joinCond string, joinType string, conds []Cond, options ...O) *joinStmt {
	return &joinStmt{
		listStmt: db.ListStmt(tblName, conds, options...),
		joinedElems: []joinedElem{
			joinedElem{
				joinType: joinType,
				joinedTbl: joinedTblName,
				joinCond: joinCond,
			},
		},
	}
}

func (db *DBI) InnerJoinStmt(tblName string, joinedTblName string, joinCond string, conds []Cond, options ...O) *joinStmt {
	return db.joinStmt(tblName, joinedTblName, joinCond, "INNER", conds, options...)
}

func (db *DBI) LeftJoinStmt(tblName string, joinedTblName string, joinCond string, conds []Cond, options ...O) *joinStmt {
	return db.joinStmt(tblName, joinedTblName, joinCond, "LEFT", conds, options...)
}

func (stmt *joinStmt) join(joinedTblName, joinCond, joinType string) *joinStmt {
	stmt.joinedElems = append(stmt.joinedElems, joinedElem{
			joinType: joinType,
			joinedTbl: joinedTblName,
			joinCond: joinCond,
	})
	return stmt
}

func (stmt *joinStmt) InnerJoin(joinedTblName string, joinCond string) *joinStmt {
	return stmt.join(joinedTblName, joinCond, "INNER")
}

func (stmt *joinStmt) LeftJoinStmt(joinedTblName string, joinCond string) *joinStmt {
	return stmt.join(joinedTblName, joinCond, "LEFT")
}

func (db *DBI) InsertStmt(tblName string, options ...O) *insertStmt {
	opts := getOptions(options...)
	return &insertStmt{
		execStmt: &execStmt{
			engine: db,
			session: opts.session,
			table: tblName,
		},
	}
}

func (db *DBI) UpdateStmt(tblName string, conds []Cond, cols []string, options ...O) *updateStmt {
	opts := getOptions(options...)
	return &updateStmt{
		execStmt: &execStmt{
			engine: db,
			session: opts.session,
			table: tblName,
			conds: conds,
		},
		cols: cols,
	}
}

func (db *DBI) UpdateSetStmt(tblName string, sets []Set, conds []Cond, options ...O) *updateSetStmt {
	opts := getOptions(options...)
	return &updateSetStmt{
		execStmt: &execStmt{
			engine: db,
			session: opts.session,
			table: tblName,
			conds: conds,
		},
		sets: sets,
	}
}

func (db *DBI) DeleteStmt(tblName string, conds []Cond, options ...O) *deleteStmt {
	opts := getOptions(options...)
	return &deleteStmt{
		execStmt: &execStmt{
			engine: db,
			session: opts.session,
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

func InnerJoinStmt(tblName string, joinedTblName string, joinCond string, conds []Cond, options ...O) *joinStmt {
	db := getDefaultConnection()
	return db.InnerJoinStmt(tblName, joinedTblName, joinCond, conds, options...)
}

func LeftJoinStmt(tblName string, joinedTblName string, joinCond string, conds []Cond, options ...O) *joinStmt {
	db := getDefaultConnection()
	return db.LeftJoinStmt(tblName, joinedTblName, joinCond, conds, options...)
}

func InsertStmt(tblName string, options ...O) *insertStmt {
	db := getDefaultConnection()
	return db.InsertStmt(tblName, options...)
}

func UpdateStmt(tblName string, conds []Cond, cols []string, options ...O) *updateStmt {
	db := getDefaultConnection()
	return db.UpdateStmt(tblName, conds, cols, options...)
}

func UpdateSetStmt(tblName string, sets []Set, conds []Cond, options ...O) *updateSetStmt {
	db := getDefaultConnection()
	return db.UpdateSetStmt(tblName, sets, conds, options...)
}

func DeleteStmt(tblName string, conds []Cond, options ...O) *deleteStmt {
	db := getDefaultConnection()
	return db.DeleteStmt(tblName, conds, options...)
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

func (db *DBI) join(tblName string, joinedTblName string, joinCond string, joinType string, conds []Cond, res interface{}, options ...O) error {
	stmt := db.joinStmt(tblName, joinedTblName, joinCond, joinType, conds, options...)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) InnerJoin(tblName string, joinedTblName string, joinCond string, conds []Cond, res interface{}, options ...O) error {
	stmt := db.InnerJoinStmt(tblName, joinedTblName, joinCond, conds, options...)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) LeftJoin(tblName string, joinedTblName string, joinCond string, conds []Cond, res interface{}, options ...O) error {
	stmt := db.LeftJoinStmt(tblName, joinedTblName, joinCond, conds, options...)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Select(tblName string, fields []string, conds []Cond, res interface{}, options ...O) error {
	stmt := db.SelectStmt(tblName, fields, conds, options...)
	_, err := stmt.Exec(res)
	return err
}

func (db *DBI) Insert(tblName string, vals interface{}, options ...O) error {
	_, err := db.InsertStmt(tblName, options...).Exec(vals)
	return err
}

func (db *DBI) Update(tblName string, conds []Cond, cols []string, vals interface{}, options ...O) (int64, error) {
	ac, err := db.UpdateStmt(tblName, conds, cols, options...).Exec(vals)
	return ac.(int64), err
}

func (db *DBI) UpdateSet(tblName string, sets []Set, conds []Cond, options ...O) (int64, error) {
	ac, err := db.UpdateSetStmt(tblName, sets, conds, options...).Exec(nil)
	return ac.(int64), err
}

func (db *DBI) Delete(tblName string, conds []Cond, vals interface{}, options ...O) error {
	_, err := db.DeleteStmt(tblName, conds, options...).Exec(vals)
	return err
}

func (db *DBI) RunSQL(tblName string, sql string, res interface{}, options ...O) error {
	stmt := db.SqlStmt(tblName, sql, options...)
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

func LeftJoin(tblName string, joinedTblName string, joinCond string, conds []Cond, res interface{}, options ...O) error {
	db := getDefaultConnection()
	return db.LeftJoin(tblName, joinedTblName, joinCond, conds, res, options...)
}

func Select(tblName string, fields []string, conds []Cond, res interface{}, options ...O) error {
	db := getDefaultConnection()
	return db.Select(tblName, fields, conds, res, options...)
}

func Insert(tblName string, vals interface{}, options ...O) error {
	db := getDefaultConnection()
	return db.Insert(tblName, vals, options...)
}

func Update(tblName string, conds []Cond, cols []string, vals interface{}, options ...O) (int64, error) {
	db := getDefaultConnection()
	return db.Update(tblName, conds, cols, vals, options...)
}

func UpdateSet(tblName string, sets []Set, conds []Cond, options ...O) (int64, error) {
	db := getDefaultConnection()
	return db.UpdateSet(tblName, sets, conds, options...)
}

func Delete(tblName string, conds []Cond, vals interface{}, options ...O) error {
	db := getDefaultConnection()
	return db.Delete(tblName, conds, vals, options...)
}

func RunSQL(tblName string, sql string, res interface{}, options ...O) error {
	db := getDefaultConnection()
	return db.RunSQL(tblName, sql, res, options...)
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
func (stmt *queryStmt) Count(bean interface{}) (int64, error) {
	sess := stmt.createQuerySession()
	return sess.Count(bean)
}

func (stmt *queryStmt) Sum(bean interface{}, col string) (float64, error) {
	sess := stmt.createQuerySession()
	return sess.Sum(bean, col)
}

func (stmt *joinStmt) Count(bean interface{}) (int64, error) {
	sess := stmt.createQuerySession()
	return sess.Count(bean)
}

func (stmt *joinStmt) Sum(bean interface{}, col string) (float64, error) {
	sess := stmt.createQuerySession()
	return sess.Sum(bean, col)
}

func getOptions(options ...O) *Options {
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}
	return opts
}

func getOneFromList(stmt Stmt, res interface{}) (has bool, err error) {
	if isSlicePtr(res) {
		if _, err = stmt.Exec(res); err != nil {
			return
		}
		has = sliceLen(res) > 0
		return
	}

	r := mk1ElemSlicePtr(res)
	if _, err = stmt.Exec(r); err != nil {
		return
	}
	if has = sliceLen(r) > 0; !has {
		return
	}
	copySliceElem(r, res)
	return
}
