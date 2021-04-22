package dbx

import (
	"strings"
	"fmt"
)

// implementations of interface Cond
type andCond struct {
	field string
	val []interface{}
}
func (e *andCond) makeCond(sess *Session) *Session {
	backquote := getQuote(e.field)
	return sess.And(fmt.Sprintf("%s%s%s=?", backquote, e.field, backquote), e.val...)
}

type andxCond struct {
	cond []string
}
func (e *andxCond) makeCond(sess *Session) *Session {
	for _, c := range e.cond {
		if len(c) > 0 {
			sess = sess.And(c)
		}
	}
	return sess
}

type opCond struct {
	field string
	op  string
	val []interface{}
}
func (e *opCond) makeCond(sess *Session) *Session {
	backquote := getQuote(e.field)
	return sess.And(fmt.Sprintf("%s%s%s %s ?", backquote, e.field, backquote, e.op), e.val...)
}

type orCond struct {
	fields []string
	vals []interface{}
}
func (e *orCond) makeCond(sess *Session) *Session {
	c := len(e.vals)
	switch c {
	case 0:
		return sess
	default:
		conds := make([]string, c)
		for i:=0; i<c; i++ {
			backquote := getQuote(e.fields[i])
			conds[i] = fmt.Sprintf("%s%s%s=?", backquote, e.fields[i], backquote)
		}
		return sess.And(strings.Join(conds, " OR "), e.vals...)
	}
}

type orxCond struct {
	conds []string
}
func (e *orxCond) makeCond(sess *Session) *Session {
	c := len(e.conds)
	switch c {
	case 0:
		return sess
	default:
		return sess.And(strings.Join(e.conds, " OR "))
	}
}

type inCond struct {
	field string
	val []interface{}
}
func (i *inCond) makeCond(sess *Session) *Session {
	return sess.In(i.field, i.val...)
}

type notInCond struct {
	Field string
	Val []interface{}
}
func (i *notInCond) makeCond(sess *Session) *Session {
	return sess.NotIn(i.Field, i.Val...)
}

type sqlCond struct {
	sql string
}
func (s *sqlCond) makeCond(sess *Session) *Session {
	return sess.Sql(s.sql)
}

func getQuote(fieldName string) (backquote string) {
	if strings.Index(fieldName, ".") < 0 {
		backquote = "`"
	}
	return
}

// implementation of interface By
type ascOrderBy struct {
	fields []string
}
func (o *ascOrderBy) makeBy(sess *Session) *Session {
	return sess.Asc(o.fields...)
}

type descOrderBy struct {
	fields []string
}
func (o *descOrderBy) makeBy(sess *Session) *Session {
	return sess.Desc(o.fields...)
}

type groupBy struct {
	field string
}
func (o *groupBy) makeBy(sess *Session) *Session {
	return sess.GroupBy(o.field)
}

// implementation of interface Limit
type limitOffset struct {
	offset int
	count int
}
func (l *limitOffset) makeLimit(sess *Session) *Session {
	if l.count > 0 {
		sess = sess.Limit(l.count, l.offset)
	}
	return sess
}
