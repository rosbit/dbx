package dbx

type dbxStmt struct {
	engine *DBI
	table string
	conds []Cond
	cols []string
	joinedTbl string
	joinCond string
	opts []O
}

func XStmt() *dbxStmt {
	db := getDefaultConnection()
	return db.XStmt()
}

func (db *DBI) XStmt() *dbxStmt {
	return &dbxStmt{
		engine: db,
	}
}

func (s *dbxStmt) Table(tbl string) *dbxStmt {
	s.table = tbl
	return s
}

func (s *dbxStmt) InnerJoin(tblName string, joinedTblName string, joinCond string) *dbxStmt {
	s.table = tblName
	s.joinedTbl = joinedTblName
	s.joinCond = joinCond
	return s
}

func (s *dbxStmt) Where(cond ...Cond) *dbxStmt {
	if len(cond) > 0 {
		s.conds = append(s.conds, cond...)
	}
	return s
}

func (s *dbxStmt) And(cond ...AndElem) *dbxStmt {
	if len(cond) > 0 {
		s.conds = append(s.conds, And(cond...))
	}
	return s
}

func (s *dbxStmt) Or(cond ...AndElem) *dbxStmt {
	if len(cond) > 0 {
		s.conds = append(s.conds, Or(cond...))
	}
	return s
}

func (s *dbxStmt) Not(cond ...AndElem) *dbxStmt {
	if len(cond) > 0 {
		s.conds = append(s.conds, Not(cond...))
	}
	return s
}

func (s *dbxStmt) In(field string, val ...interface{}) *dbxStmt {
	s.conds = append(s.conds, In(field, val...))
	return s
}

func (s *dbxStmt) NotIn(field string, val ...interface{}) *dbxStmt {
	s.conds = append(s.conds, NotIn(field, val...))
	return s
}

func (s *dbxStmt) Cols(col ...string) *dbxStmt {
	if len(col) > 0 {
		s.cols = append(s.cols, col...)
	}
	return s
}

func (s *dbxStmt) Desc(field ...string) *dbxStmt {
	if len(field) > 0 {
		s.opts = append(s.opts, OrderByDesc(field...))
	}
	return s
}

func (s *dbxStmt) Asc(field ...string) *dbxStmt {
	if len(field) > 0 {
		s.opts = append(s.opts, OrderByAsc(field...))
	}
	return s
}

func (s *dbxStmt) GroupBy(field ...string) *dbxStmt {
	if len(field) > 0 {
		s.opts = append(s.opts, GroupBy(field...))
	}
	return s
}

func (s *dbxStmt) Limit(count int, offset ...int) *dbxStmt {
	s.opts = append(s.opts, Limit(count, offset...))
	return s
}

func (s *dbxStmt) XSession(session *Session) *dbxStmt {
	s.opts = append(s.opts, WithSession(session))
	return s
}

func (s *dbxStmt) Get(res interface{}) (bool, error) {
	return s.engine.Get(s.table, s.conds, res, s.opts...)
}

func (s *dbxStmt) List(res interface{}) error {
	if len(s.joinedTbl) > 0 && len(s.joinCond) > 0 {
		return s.engine.InnerJoin(s.table, s.joinedTbl, s.joinCond, s.conds, res, s.opts...)
	}
	return s.engine.List(s.table, s.conds, res, s.opts...)
}

func (s *dbxStmt) Insert(vals interface{}) error {
	return s.engine.Insert(s.table, vals, s.opts...)
}

func (s *dbxStmt) Update(vals interface{}) error {
	return s.engine.Update(s.table, s.conds, s.cols, vals, s.opts...)
}

func (s *dbxStmt) Delete(vals interface{}) error {
	return s.engine.Delete(s.table, s.conds, vals, s.opts...)
}

func (s *dbxStmt) Iter(bean interface{}) (<-chan interface{}) {
	return s.engine.Iter(s.table, s.conds, bean, s.opts...)
}

func (s *dbxStmt) Iterate(bean interface{}, it FnIterate) error {
	return s.engine.Iterate(s.table, s.conds, bean, it, s.opts...)
}

func (s *dbxStmt) Count(bean interface{}) (int64, error) {
	if len(s.joinedTbl) > 0 && len(s.joinCond) > 0 {
		return s.engine.InnerJoinStmt(s.table, s.joinedTbl, s.joinCond, s.conds, s.opts...).Count(bean)
	}
	return s.engine.ListStmt(s.table, s.conds, s.opts...).Count(bean)
}

func (s *dbxStmt) Sum(bean interface{}, col string) (float64, error) {
	if len(s.joinedTbl) > 0 && len(s.joinCond) > 0 {
		return s.engine.InnerJoinStmt(s.table, s.joinedTbl, s.joinCond, s.conds, s.opts...).Sum(bean, col)
	}
	return s.engine.ListStmt(s.table, s.conds, s.opts...).Sum(bean, col)
}
