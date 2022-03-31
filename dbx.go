package dbx

type dbxStmt struct {
	engine *DBI
	table string
	conds []Cond
	sets []Set
	cols []string
	joinedElems []joinedElem
	opts []O
	selection string
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
		s.joinedElems = nil
		s.opts = nil
		s.selection = ""
	}
	return s
}

func (s *dbxStmt) join(tblName string, joinedTblName string, joinCond string, joinType string) *dbxStmt {
	s.table = tblName
	s.joinedElems = []joinedElem{
		joinedElem{
			joinedTbl: joinedTblName,
			joinCond: joinCond,
			joinType: joinType,
		},
	}
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

func (s *dbxStmt) nextJoin(joinedTblName string, joinCond string, joinType string) *dbxStmt {
	s.joinedElems = append(s.joinedElems, joinedElem{
		joinedTbl: joinedTblName,
		joinCond: joinCond,
		joinType: joinType,
	})
	return s
}

func (s *dbxStmt) NextInnerJoin(joinedTblName string, joinCond string) *dbxStmt {
	s.nextJoin(joinedTblName, joinCond, "INNER")
	return s
}

func (s *dbxStmt) NextLeftJoin(joinedTblName string, joinCond string) *dbxStmt {
	s.nextJoin(joinedTblName, joinCond, "LEFT")
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
		for _, c := range cond {
			s.conds = append(s.conds, c)
		}
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

// called after InnerJoin/LeftJoin
func (s *dbxStmt) SelectCols(selection string) *dbxStmt {
	if len(selection) > 0 {
		s.opts = append(s.opts, SelectCols(selection))
		s.selection = selection
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
	if len(s.joinedElems) > 0 {
		s.opts = append(s.opts, Limit(1))
		stmt := s.generateJoinStmt()
		return getOneFromList(stmt, res)
	} else if len(s.selection) > 0 {
		s.opts = append(s.opts, Limit(1))
		stmt := s.engine.SelectStmt(s.table, []string{s.selection}, s.conds, s.opts...)
		return getOneFromList(stmt, res)
	}
	return s.engine.Get(s.table, s.conds, res, s.opts...)
}

func (s *dbxStmt) List(res interface{}) error {
	if stmt := s.generateJoinStmt(); stmt != nil {
		_, err := stmt.Exec(res)
		return err
	} else if len(s.selection) > 0 {
		return s.engine.Select(s.table, []string{s.selection}, s.conds, res, s.opts...)
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
	if stmt := s.generateJoinStmt(); stmt != nil {
		return stmt.Iter(bean)
	}
	return s.engine.Iter(s.table, s.conds, bean, s.opts...)
}

func (s *dbxStmt) Iterate(bean interface{}, it FnIterate) error {
	if stmt := s.generateJoinStmt(); stmt != nil {
		return stmt.Iterate(bean, it)
	}
	return s.engine.Iterate(s.table, s.conds, bean, it, s.opts...)
}

func (s *dbxStmt) Count(bean interface{}) (int64, error) {
	if stmt := s.generateJoinStmt(); stmt != nil {
		return stmt.Count(bean)
	}
	return s.engine.ListStmt(s.table, s.conds, s.opts...).Count(bean)
}

func (s *dbxStmt) Sum(bean interface{}, col string) (float64, error) {
	if stmt := s.generateJoinStmt(); stmt != nil {
		return stmt.Sum(bean, col)
	}
	return s.engine.ListStmt(s.table, s.conds, s.opts...).Sum(bean, col)
}

func (s *dbxStmt) generateJoinStmt() *joinStmt {
	if len(s.joinedElems) == 0 {
		return nil
	}
	e0 := &s.joinedElems[0]
	stmt := s.engine.joinStmt(s.table, e0.joinedTbl, e0.joinCond, e0.joinType, s.conds, s.opts...)
	for i:=1; i<len(s.joinedElems); i++ {
		e := &s.joinedElems[i]
		stmt.join(e.joinedTbl, e.joinCond, e.joinType)
	}
	return stmt
}
