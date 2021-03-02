package dbx

import (
	"strings"
	"fmt"
)

// implementations of interface Eq
type andCond struct {
	field string
	val []interface{}
}
func (e *andCond) makeCond(sess *Session) *Session {
	backquote := getQuote(e.field)
	return sess.And(fmt.Sprintf("%s%s%s=?", backquote, e.field, backquote), e.val...)
}

type opCond struct {
	field string
	op  string
	val interface{}
}
func (e *opCond) makeCond(sess *Session) *Session {
	backquote := getQuote(e.field)
	return sess.And(fmt.Sprintf("%s%s%s %s ?", backquote, e.field, backquote, e.op), e.val)
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

