package dbx

type dbxStmt struct {
	engine *DBI
	table string
	conds []Cond
	sets []Set
	cols []string
	joinedTbl string
	joinCond string
	joinType string
	opts []O
}

func XStmt(tbl ...string) *dbxStmt {
	db := getDefaultConnection()
	return db.XStmt(tbl...)
}

func (db *DBI) XStmt(tbl ...string) *dbxStmt {
	var t string
	if len(tbl) > 0 {
		t = tbl[0]
	}

	return &dbxStmt{
		engine: db,
		table: t,
	}
}

func (s *dbxStmt) Table(tbl string, dontReset ...bool) *dbxStmt {
	s.table = tbl
	if !(len(dontReset) > 0 && dontReset[0]) {
		s.conds = nil
		s.sets = nil
		s.cols = nil
		s.joinedTbl = ""
		s.joinCond = ""
		s.joinType = ""
		s.opts = nil
	}
	return s
}

func (s *dbxStmt) join(tblName string, joinedTblName string, joinCond string, joinType string) *dbxStmt {
	s.table = tblName
	s.joinedTbl = joinedTblName
	s.joinCond = joinCond
	s.joinType = joinType
	return s
}

func (s *dbxStmt) InnerJoin(tblName string, joinedTblName string, joinCond string) *dbxStmt {
	s.join(tblName, joinedTblName, joinCond, "INNER")
	return s
}

func (s *dbxStmt) LeftJoin(tblName string, joinedTblName string, joinCond string) *dbxStmt {
	s.join(tblName, joinedTblName, joinCond, "LEFT")
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

func (s *dbxStmt) Set(set ...Set) *dbxStmt {
	if len(set) > 0 {
		s.sets = append(s.sets, set...)
	}
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

func (s *dbxStmt) SelectCols(selection string) *dbxStmt {
	if len(selection) > 0 {
		s.opts = append(s.opts, SelectCols(selection))
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

func (s *dbxStmt) Get(res interface{}) (has bool, err error) {
	if len(s.joinedTbl) > 0 && len(s.joinCond) > 0 {
		s.opts = append(s.opts, Limit(1))
		if isSlicePtr(res) {
			if err = s.engine.join(s.table, s.joinedTbl, s.joinCond, s.joinType, s.conds, res, s.opts...); err != nil {
				return
			}
			has = sliceLen(res) > 0
			return
		}

		r := mk1ElemSlicePtr(res)
		if err = s.engine.join(s.table, s.joinedTbl, s.joinCond, s.joinType, s.conds, r, s.opts...); err != nil {
			return
		}
		if has = sliceLen(r) > 0; !has {
			return
		}
		copySliceElem(r, res)
		return
	}
	return s.engine.Get(s.table, s.conds, res, s.opts...)
}

func (s *dbxStmt) List(res interface{}) error {
	if len(s.joinedTbl) > 0 && len(s.joinCond) > 0 {
		return s.engine.join(s.table, s.joinedTbl, s.joinCond, s.joinType, s.conds, res, s.opts...)
	}
	return s.engine.List(s.table, s.conds, res, s.opts...)
}

func (s *dbxStmt) Insert(vals interface{}) error {
	return s.engine.Insert(s.table, vals, s.opts...)
}

func (s *dbxStmt) Update(vals interface{}) (int64, error) {
	if len(s.sets) == 0 {
		return s.engine.Update(s.table, s.conds, s.cols, vals, s.opts...)
	}
	return s.engine.UpdateSet(s.table, s.sets, s.conds, s.opts...)
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
		return s.engine.joinStmt(s.table, s.joinedTbl, s.joinCond, s.joinType, s.conds, s.opts...).Count(bean)
	}
	return s.engine.ListStmt(s.table, s.conds, s.opts...).Count(bean)
}

func (s *dbxStmt) Sum(bean interface{}, col string) (float64, error) {
	if len(s.joinedTbl) > 0 && len(s.joinCond) > 0 {
		return s.engine.joinStmt(s.table, s.joinedTbl, s.joinCond, s.joinType, s.conds, s.opts...).Sum(bean, col)
	}
	return s.engine.ListStmt(s.table, s.conds, s.opts...).Sum(bean, col)
}
