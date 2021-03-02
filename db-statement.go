package dbx

import (
	"strings"
	"fmt"
)

const (
	_select = "SELECT"
	_sql = "SQL"
	_innerJoin = "INNER"
)

type execStmt struct {
	engine *DBI
	table   string
	eqs   []Eq
	conds []string
}

func (stmt *execStmt) createExecSession(session *Session, extraQuery ...map[string]interface{}) *Session {
	var sess *Session
	if session == nil {
		sess = stmt.engine.Table(stmt.table)
	} else {
		sess = session.Table(stmt.table)
	}
	if len(extraQuery) > 0 {
		for k, v := range extraQuery[0] {
			switch k {
			case _select:
				if fields, ok := v.([]string); ok {
					sess.Select(strings.Join(fields, ","))
				}
			case _sql:
				if sql, ok := v.(string); ok {
					sess.SQL(sql)
				}
			case _innerJoin:
				vals := v.([]string)
				joinedTbl, joinCond := vals[0], vals[1]
				sess = sess.Select(fmt.Sprintf("%s.*, %s.*", stmt.table, joinedTbl)).
					Join("INNER", joinedTbl, joinCond)
			default:
			}
		}
	}

	hasWhere := false
	if len(stmt.eqs) > 0 {
		hasWhere = true
		for i, _ := range stmt.eqs {
			sess = stmt.eqs[i].makeCond(sess)
		}
	}
	if len(stmt.conds) > 0 {
		startIdx := 0
		if !hasWhere {
			hasWhere = true
			sess = sess.Where(stmt.conds[0])
			startIdx = 1
		}
		for i:=startIdx; i<len(stmt.conds); i++ {
			sess = sess.And(stmt.conds[i])
		}
	}

	return sess
}

type queryStmt struct {
	*execStmt
	sortDesc []string
	sortAsc  []string
}
func (stmt *queryStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	sess := stmt.createQuerySession(session)
	return sess.Get(bean)
}

func (stmt *queryStmt) createQuerySession(session *Session, extraQuery ...map[string]interface{}) *Session {
	sess := stmt.execStmt.createExecSession(session, extraQuery...)

	if len(stmt.sortDesc) > 0 {
		sess = sess.Desc(stmt.sortDesc...)
	}
	if len(stmt.sortAsc) > 0 {
		sess = sess.Asc(stmt.sortAsc...)
	}

	return sess
}

type limitT struct {
	offset int
	count int
}
type listStmt struct {
	*queryStmt
	limit limitT
}
func (stmt *listStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	sess := stmt.queryStmt.createQuerySession(session)
	return stmt.find(sess, bean)
}

func (stmt *listStmt) find(sess *Session, bean interface{}) (StmtResult, error) {
	if stmt.limit.count > 0 {
		sess = sess.Limit(stmt.limit.count, stmt.limit.offset)
	}
	err := sess.Find(bean)
	return nil, err
}

type selectStmt struct {
	*listStmt
	fields []string
}
func (stmt *selectStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	sess := stmt.queryStmt.createQuerySession(session, map[string]interface{}{_select:stmt.fields})
	return stmt.find(sess, bean)
}

type sqlStmt struct {
	*listStmt
	sql string
}
func (stmt *sqlStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	sess := stmt.queryStmt.createQuerySession(session, map[string]interface{}{_sql:stmt.sql})
	return stmt.find(sess, bean)
}

type innerJoinStmt struct {
	*listStmt
	joinedTbl string
	joinCond string
}
func (stmt *innerJoinStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	sess := stmt.queryStmt.createQuerySession(session, map[string]interface{}{
		_innerJoin:[]string{stmt.joinedTbl, stmt.joinCond},
	})
	return stmt.find(sess, bean)
}

type updateStmt struct {
	*execStmt
	cols []string
}
func (stmt *updateStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	sess := stmt.execStmt.createExecSession(session)
	if len(stmt.cols) > 0 {
		sess = sess.Cols(stmt.cols...)
	}
	return sess.Update(bean)
}

type insertStmt struct {
	*execStmt
}
func (stmt *insertStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	stmt.eqs, stmt.conds = nil, nil
	sess := stmt.execStmt.createExecSession(session)
	return sess.Insert(bean)
}

type deleteStmt struct {
	*execStmt
}
func (stmt *deleteStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	sess := stmt.execStmt.createExecSession(session)
	return sess.Delete(bean)
}

type voidStmt struct {
}

func (stmt *voidStmt) Exec(bean interface{}, session *Session) (StmtResult, error) {
	return true, nil
}

