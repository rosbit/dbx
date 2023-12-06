package dbx

import (
	"database/sql"
	"strings"
	"fmt"
)

const (
	_select = "SELECT"
	_sql = "SQL"
	_join = "JOIN"
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
			case _join:
				jStmt := v.(*joinStmt)
				if len(jStmt.selection) == 0 {
					tbls := make([]string, len(jStmt.joinedElems)+1)
					tbls[0] = fmt.Sprintf("%s.*", stmt.table)
					for i, e := range jStmt.joinedElems {
						tbls[i+1] = fmt.Sprintf("%s.*", e.joinedTbl)
					}
					sess = sess.Select(strings.Join(tbls, ","))
				} else {
					sess = sess.Select(jStmt.selection)
				}
				for _, e := range jStmt.joinedElems {
					sess.Join(e.joinType, e.joinedTbl, e.joinCond)
				}
			default:
			}
		}
	}

	if len(stmt.conds) > 0 {
		sess1 := (*xormSession)(sess)
		buildConds(sess1, stmt.conds)
		sess = (*Session)(sess1)
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
	selection string
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

type joinedElem struct {
	joinType string
	joinedTbl string
	joinCond string
}
type joinStmt struct {
	*listStmt
	joinedElems []joinedElem
}
func (stmt *joinStmt) createQuerySession() *Session {
	return stmt.queryStmt.createQuerySession(map[string]interface{}{
		_join: stmt,
	})
}
func (stmt *joinStmt) Exec(bean interface{}) (StmtResult, error) {
	sess := stmt.createQuerySession()
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

type rawUpdateStmt struct {
	*execStmt
	sql string
}
func (stmt *rawUpdateStmt) Exec(_ interface{}) (StmtResult, error) {
	sess := stmt.execStmt.createExecSession()
	return sess.Exec(stmt.sql)
}

type updateSetStmt struct {
	*execStmt
	sets []Set
}
func (stmt *updateSetStmt) Exec(bean interface{}) (StmtResult, error) {
	if len(stmt.sets) == 0 || len(stmt.table) == 0 {
		return int64(0), nil
	}

	sb := newSqlBuilder()
	fmt.Fprintf(sb.q, "UPDATE %s SET", stmt.table)
	sb.appendSets(stmt.sets)
	if len(stmt.conds) > 0 {
		sb.q.WriteString(" WHERE")
		buildConds(sb, stmt.conds)
	}
	params := sb.toParams()

	var sess *Session
	if stmt.session == nil {
		sess = stmt.engine.Table(stmt.table)
	} else {
		sess = stmt.session.Table(stmt.table)
	}
	r, err := sess.Exec(params...)
	if err != nil {
		return int64(0), err
	}
	if r1, ok := r.(sql.Result); ok {
		return r1.RowsAffected()
	}
	return int64(0), nil
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
