package dbx

import (
	"strings"
	"fmt"
)

func makeCond(a AndElem, sess *Session) *Session {
    q, v := a.mkAndElem()
    if len(q) > 0 {
        sess = sess.And(q, v...)
    }
    return sess
}

// join with And/OR
func joinAndElems(conds []AndElem, conj string) (string, []interface{}) {
	if len(conds) == 0 {
		return "", nil
	}

	var elemFmt string
	if len(conds) == 1 {
		elemFmt = "%s"
	} else {
		elemFmt = "(%s)"
	}

	res := &strings.Builder{}
	vals := []interface{}{}
	for _, c := range conds {
		if c != nil {
			q, v := c.mkAndElem()

			if len(q) > 0 {
				if res.Len() > 0 {
					fmt.Fprintf(res, " %s ", conj)
				}
				fmt.Fprintf(res, elemFmt, q)

				if len(v) > 0 {
					vals = append(vals, v...)
				}
			}
		}
	}
	return res.String(), vals
}

// make "IN" or "NOT IN"
func makeInElem(field string, val []interface{}, prep string) (string, []interface{}) {
	if len(field) == 0 || len(val) == 0 {
		return "", nil
	}

	res := &strings.Builder{}
	backquote := getQuote(field)
	fmt.Fprintf(res, "%s%s%s %s ", backquote, field, backquote, prep)
	for i, _ := range val {
		if i == 0 {
			fmt.Fprintf(res, "(?")
		} else {
			fmt.Fprintf(res, ",?")
		}
	}
	fmt.Fprintf(res, ")")
	return res.String(), val
}

// implementations of interface Cond
type andCond struct {
	field string
	val interface{}
}
func (e *andCond) makeCond(sess *Session) *Session {
	return makeCond(e, sess)
}

func (e *andCond) mkAndElem() (string, []interface{}) {
	if len(e.field) == 0 {
		return "", nil
	}
	backquote := getQuote(e.field)
	return fmt.Sprintf("%s%s%s=?", backquote, e.field, backquote), []interface{}{e.val}
}

type andxCond struct {
	conds []AndElem
}
func (e *andxCond) makeCond(sess *Session) *Session {
	return makeCond(e, sess)
}
func (e *andxCond) mkAndElem() (string, []interface{}) {
	return joinAndElems(e.conds, "AND")
}

type opCond struct {
	field string
	op  string
	val interface{}
}
func (e *opCond) makeCond(sess *Session) *Session {
	return makeCond(e, sess)
}
func (e *opCond) mkAndElem() (string, []interface{}) {
	if len(e.field) == 0 {
		return "", nil
	}
	if len(e.op) == 0 {
		e.op = "="
	}
	backquote := getQuote(e.field)
	return fmt.Sprintf("%s%s%s %s ?", backquote, e.field, backquote, e.op), []interface{}{e.val}
}

type orxCond struct {
	conds []AndElem
}
func (e *orxCond) makeCond(sess *Session) *Session {
	return makeCond(e, sess)
}
func (e *orxCond) mkAndElem() (string, []interface{}) {
	return joinAndElems(e.conds, "OR")
}

type notCond struct {
	conds []AndElem
}
func (e *notCond) makeCond(sess *Session) *Session {
	return makeCond(e, sess)
}
func (e *notCond) mkAndElem() (string, []interface{}) {
	q, v := joinAndElems(e.conds, "NOT")
	if len(q) == 0 {
		return q, v
	}
	return fmt.Sprintf("NOT (%s)", q), v
}

type inCond struct {
	field string
	val []interface{}
}
func (i *inCond) makeCond(sess *Session) *Session {
	return makeCond(i, sess)
}
func (i *inCond) mkAndElem() (string, []interface{}) {
	return makeInElem(i.field, i.val, "IN")
}

type notInCond struct {
	field string
	val []interface{}
}
func (i *notInCond) makeCond(sess *Session) *Session {
	return makeCond(i, sess)
}
func (i *notInCond) mkAndElem() (string, []interface{}) {
	return makeInElem(i.field, i.val, "NOT IN")
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
