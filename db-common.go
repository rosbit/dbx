package dbx

import (
	"fmt"
)

func (db *DBI) NewQueryStmt(tblName string, params []Eq, conds []string, descFields []string, ascFields []string) *queryStmt {
	return &queryStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			eqs: params,
			conds: conds,
		},
		sortDesc: descFields,
		sortAsc: ascFields,
	}
}

func (db *DBI) NewListStmt(tblName string, params []Eq, conds []string, descFields []string, ascFields []string, offset, count int) *listStmt {
	return &listStmt{
		queryStmt: db.NewQueryStmt(tblName, params, conds, descFields, ascFields),
		limit: limitT{offset:offset, count:count},
	}
}

func (db *DBI) NewSelectStmt(tblName string, fields []string, params []Eq, conds []string, descFields []string, ascFields []string, offset, count int) *selectStmt {
	return &selectStmt{
		listStmt: db.NewListStmt(tblName, params, conds, descFields, ascFields, offset, count),
		fields: fields,
	}
}

func (db *DBI) NewSqlStmt(tblName string, sql string, params []Eq, conds []string, descFields []string, ascFields []string, offset, count int) *sqlStmt {
	return &sqlStmt{
		listStmt: db.NewListStmt(tblName, params, conds, descFields, ascFields, offset, count),
		sql: sql,
	}
}

func (db *DBI) NewInnerJoinStmt(tblName string, joinedTblName string, joinCond string, params []Eq, conds []string, descFields []string, ascFields []string, offset, count int) *innerJoinStmt {
	return &innerJoinStmt{
		listStmt: db.NewListStmt(tblName, params, conds, descFields, ascFields, offset, count),
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

func (db *DBI) NewUpdateStmt(tblName string, params []Eq, conds []string, cols []string) *updateStmt {
	return &updateStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			eqs: params,
			conds: conds,
		},
		cols: cols,
	}
}

func (db *DBI) NewDeleteStmt(tblName string, params []Eq, conds []string) *deleteStmt {
	return &deleteStmt{
		execStmt: &execStmt{
			engine: db,
			table: tblName,
			eqs: params,
			conds: conds,
		},
	}
}

func NewQueryStmt(tblName string, params []Eq, conds []string, descFields []string, ascFields []string) *queryStmt {
	db := getDefaultConnection()
	return db.NewQueryStmt(tblName, params, conds, descFields, ascFields)
}

func NewListStmt(tblName string, params []Eq, conds []string, descFields []string, ascFields []string, offset, count int) *listStmt {
	db := getDefaultConnection()
	return db.NewListStmt(tblName, params, conds, descFields, ascFields, offset, count)
}

func NewSelectStmt(tblName string, fields []string, params []Eq, conds []string, descFields []string, ascFields []string, offset, count int) *selectStmt {
	db := getDefaultConnection()
	return db.NewSelectStmt(tblName, fields, params, conds, descFields, ascFields, offset, count)
}

func NewSqlStmt(tblName string, sql string, params []Eq, conds []string, descFields []string, ascFields []string, offset, count int) *sqlStmt {
	db := getDefaultConnection()
	return db.NewSqlStmt(tblName, sql, params, conds, descFields, ascFields, offset, count)
}

func NewInnerJoinStmt(tblName string, joinedTblName string, joinCond string, params []Eq, conds []string, descFields []string, ascFields []string, offset, count int) *innerJoinStmt {
	db := getDefaultConnection()
	return db.NewInnerJoinStmt(tblName, joinedTblName, joinCond, params, conds, descFields, ascFields, offset, count)
}

func NewInsertStmt(tblName string) *insertStmt {
	db := getDefaultConnection()
	return db.NewInsertStmt(tblName)
}

func NewUpdateStmt(tblName string, params []Eq, conds []string, cols []string) *updateStmt {
	db := getDefaultConnection()
	return db.NewUpdateStmt(tblName, params, conds, cols)
}

func NewDeleteStmt(tblName string, params []Eq, conds []string) *deleteStmt {
	db := getDefaultConnection()
	return db.NewDeleteStmt(tblName, params, conds)
}

func NewVoidStmt() *voidStmt {
	return &voidStmt{}
}

func NewEq(fieldName string, val ...interface{}) Eq {
	return &andCond{fieldName, val}
}

// op: "=", "!=", "<>", ">", ">=", "<", "<=", "like"
func NewOp(fieldName string, op string, val interface{}) Eq {
	return &opCond{fieldName, op, val}
}

// f1, v1, f2, v2, ... => (f1=v1 OR f2=v2 OR ...)
func NewOr(fieldName string, val ...interface{}) Eq {
	fields := []string{fieldName}
	vals := []interface{}{}

	for i, v := range val {
		if i % 2 == 0 {
			vals = append(vals, v)
		} else {
			fields = append(fields, fmt.Sprintf("%s", v))
		}
	}
	return &orCond{fields, vals}
}

func NewIn(fieldName string, val ...interface{}) Eq {
	return &inCond{fieldName, copy_i(val)}
}

func NewNotIn(fieldName string, val ...interface{}) Eq {
	return &notInCond{fieldName, copy_i(val)}
}

func copy_i(vals ...interface{}) []interface{} {
	if len(vals) == 0 {
		return nil
	}
	vs := make([]interface{}, len(vals))
	for i, _ := range vals {
		vs[i] = vals[i]
	}
	return vs
}

func NewSql(sql string) Eq {
	return &sqlCond{sql}
}

// some re-usable handler
func (db *DBI) GetById(tblName, idName string, idVal interface{}, res interface{}) (bool, error) {
	eqs := []Eq{NewEq(idName, idVal)}
	stmt := db.NewQueryStmt(tblName, eqs, nil, nil, nil)
	res, err := stmt.Exec(res, nil)
	if err != nil {
		return false, err
	}
	return res.(bool), err
}

func (db *DBI) GetOne(tblName string, params []Eq, conds []string, res interface{}) (bool, error) {
	stmt := db.NewQueryStmt(tblName, params, conds, nil, nil)
	res, err := stmt.Exec(res, nil)
	if err != nil {
		return false, err
	}
	return res.(bool), err
}

func (db *DBI) Find(tblName string, params []Eq, conds []string, offset, count int, res interface{}) error {
	stmt := db.NewListStmt(tblName, params, conds, nil, nil, offset, count)
	_, err := stmt.Exec(res, nil)
	return err
}

func (db *DBI) Select(tblName string, fields []string, res interface{}) error {
	stmt := db.NewSelectStmt(tblName, fields, nil, nil, nil, nil, 0, 0)
	_, err := stmt.Exec(res, nil)
	return err
}

func (db *DBI) SQL(tblName string, sql string, res interface{}) error {
	stmt := db.NewSqlStmt(tblName, sql, nil, nil, nil, nil, 0, 0)
	_, err := stmt.Exec(res, nil)
	return err
}

func (db *DBI) Iter(tblName string, params []Eq, conds []string, bean interface{}) (<-chan interface{}) {
	stmt := db.NewQueryStmt(tblName, params, conds, nil, nil)
	return stmt.Iter(bean)
}

func (db *DBI) Iterate(tblName string, params []Eq, conds []string, bean interface{}, it FnIterate) error {
	stmt := db.NewQueryStmt(tblName, params, conds, nil, nil)
	return stmt.Iterate(bean, it)
}

func GetById(tblName, idName string, idVal interface{}, res interface{}) (bool, error) {
	db := getDefaultConnection()
	return db.GetById(tblName, idName, idVal, res)
}

func GetOne(tblName string, params []Eq, conds []string, res interface{}) (bool, error) {
	db := getDefaultConnection()
	return db.GetOne(tblName, params, conds, res)
}

func Find(tblName string, params []Eq, conds []string, offset, count int, res interface{}) error {
	db := getDefaultConnection()
	return db.Find(tblName, params, conds, offset, count, res)
}

func Select(tblName string, fields []string, res interface{}) error {
	db := getDefaultConnection()
	return db.Select(tblName, fields, res)
}

func SQL(tblName string, sql string, res interface{}) error {
	db := getDefaultConnection()
	return db.SQL(tblName, sql, res)
}

func Iter(tblName string, params []Eq, conds []string, bean interface{}) (<-chan interface{}) {
	db := getDefaultConnection()
	return db.Iter(tblName, params, conds, bean)
}

func Iterate(tblName string, params []Eq, conds []string, bean interface{}, it FnIterate) error {
	db := getDefaultConnection()
	return db.Iterate(tblName, params, conds, bean, it)
}

// some statistic func
func (stmt *queryStmt) Count(bean interface{}, session *Session) (int64, error) {
	sess := stmt.createQuerySession(session)
	return sess.Count(bean)
}

