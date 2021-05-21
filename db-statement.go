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
	engine  *DBI
	session *Session
	table   string
	conds   []Cond
}

func (stmt *execStmt) createExecSession(extraQuery ...map[string]interface{}) *Session {
	var sess *Session
	if stmt.session == nil {
		sess = stmt.engine.Table(stmt.table)
	} else {
		sess = stmt.session.Table(stmt.table)
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

	if len(stmt.conds) > 0 {
		for i, _ := range stmt.conds {
			sess = stmt.conds[i].makeCond(sess)
		}
	}

	return sess
}

func (stmt *execStmt) setSession(session *Session) {
	if stmt.session == nil {
		stmt.session = session
	}
}

type queryStmt struct {
	*execStmt
	bys []by
	limit limit
}
func (stmt *queryStmt) Exec(bean interface{}) (StmtResult, error) {
	sess := stmt.createQuerySession()
	return sess.Get(bean)
}

func (stmt *queryStmt) createQuerySession(extraQuery ...map[string]interface{}) *Session {
	sess := stmt.execStmt.createExecSession(extraQuery...)

	for _, b := range stmt.bys {
		sess = b.makeBy(sess)
	}

	if stmt.limit != nil {
		sess = stmt.limit.makeLimit(sess)
	}
	return sess
}

type listStmt struct {
	*queryStmt
}
func (stmt *listStmt) Exec(bean interface{}) (StmtResult, error) {
	sess := stmt.queryStmt.createQuerySession()
	return stmt.find(sess, bean)
}

func (stmt *listStmt) find(sess *Session, bean interface{}) (StmtResult, error) {
	err := sess.Find(bean)
	return nil, err
}

type selectStmt struct {
	*listStmt
	fields []string
}
func (stmt *selectStmt) Exec(bean interface{}) (StmtResult, error) {
	sess := stmt.queryStmt.createQuerySession(map[string]interface{}{_select:stmt.fields})
	return stmt.find(sess, bean)
}

type sqlStmt struct {
	*listStmt
	sql string
}
func (stmt *sqlStmt) Exec(bean interface{}) (StmtResult, error) {
	sess := stmt.queryStmt.createQuerySession(map[string]interface{}{_sql:stmt.sql})
	return stmt.find(sess, bean)
}

type innerJoinStmt struct {
	*listStmt
	joinedTbl string
	joinCond string
}
func (stmt *innerJoinStmt) Exec(bean interface{}) (StmtResult, error) {
	sess := stmt.queryStmt.createQuerySession(map[string]interface{}{
		_innerJoin:[]string{stmt.joinedTbl, stmt.joinCond},
	})
	return stmt.find(sess, bean)
}

type updateStmt struct {
	*execStmt
	cols []string
}
func (stmt *updateStmt) Exec(bean interface{}) (StmtResult, error) {
	sess := stmt.execStmt.createExecSession()
	if len(stmt.cols) > 0 {
		sess = sess.Cols(stmt.cols...)
	}
	return sess.Update(bean)
}

type insertStmt struct {
	*execStmt
}
func (stmt *insertStmt) Exec(bean interface{}) (StmtResult, error) {
	stmt.conds = nil
	sess := stmt.execStmt.createExecSession()
	return sess.Insert(bean)
}

type deleteStmt struct {
	*execStmt
}
func (stmt *deleteStmt) Exec(bean interface{}) (StmtResult, error) {
	sess := stmt.execStmt.createExecSession()
	return sess.Delete(bean)
}
