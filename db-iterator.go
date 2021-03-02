package dbx

// ---- BEGIN: iterate result set with channel ----
func (stmt *queryStmt) Iter(bean interface{}) (<-chan interface{}) {
	sess := stmt.createQuerySession(nil)
	return iter(sess, bean)
}

func (stmt *listStmt) Iter(bean interface{}) (<-chan interface{}) {
	sess := stmt.queryStmt.createQuerySession(nil)
	return stmt.iter(sess, bean)
}

func (stmt *selectStmt) Iter(bean interface{}) (<-chan interface{}) {
	sess := stmt.queryStmt.createQuerySession(nil, map[string]interface{}{_select:stmt.fields})
	return stmt.listStmt.iter(sess, bean)
}

func (stmt *sqlStmt) Iter(bean interface{}) (<-chan interface{}) {
	sess := stmt.queryStmt.createQuerySession(nil, map[string]interface{}{_sql:stmt.sql})
	return stmt.listStmt.iter(sess, bean)
}

func (stmt *listStmt) iter(sess *Session, bean interface{}) (<-chan interface{}) {
	if stmt.limit.count > 0 {
		sess = sess.Limit(stmt.limit.count, stmt.limit.offset)
	}
	return iter(sess, bean)
}

func (stmt *innerJoinStmt) Iter(bean interface{}) (<-chan interface{}) {
	sess := stmt.queryStmt.createQuerySession(nil, map[string]interface{}{
		_innerJoin:[]string{stmt.joinedTbl, stmt.joinCond},
	})
	return stmt.listStmt.iter(sess, bean)
}

func iter(sess *Session, bean interface{}) (<-chan interface{}) {
	c := make(chan interface{})
	go func() {
		sess.Iterate(bean, func(_ int, bean interface{})error{
			c <- bean
			return nil
		})
		close(c)
	}()
	return c
}
// ---- END: iterate result set with channel ----

// ---- BEGIN: iterate result set using callback ----
func (stmt *queryStmt) Iterate(bean interface{}, it FnIterate) error {
	sess := stmt.createQuerySession(nil)
	return sess.Iterate(bean, it)
}

func (stmt *listStmt) Iterate(bean interface{}, it FnIterate) error {
	sess := stmt.queryStmt.createQuerySession(nil)
	return stmt.iterate(sess, bean, it)
}

func (stmt *selectStmt) Iterate(bean interface{}, it FnIterate) error {
	sess := stmt.queryStmt.createQuerySession(nil, map[string]interface{}{_select:stmt.fields})
	return stmt.listStmt.iterate(sess, bean, it)
}

func (stmt *sqlStmt) Iterate(bean interface{}, it FnIterate) error {
	sess := stmt.queryStmt.createQuerySession(nil, map[string]interface{}{_sql:stmt.sql})
	return stmt.listStmt.iterate(sess, bean, it)
}

func (stmt *innerJoinStmt) Iterate(bean interface{}, it FnIterate) error {
	sess := stmt.queryStmt.createQuerySession(nil, map[string]interface{}{
		_innerJoin:[]string{stmt.joinedTbl, stmt.joinCond},
	})
	return stmt.listStmt.iterate(sess, bean, it)
}

func (stmt *listStmt) iterate(sess *Session, bean interface{}, it FnIterate) error {
	if stmt.limit.count > 0 {
		sess = sess.Limit(stmt.limit.count, stmt.limit.offset)
	}
	return sess.Iterate(bean, it)
}

// ---- END: iterate result set using callback ----

