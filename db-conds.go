package dbx

import (
	"strings"
	"fmt"
)

type andElemWrapper struct {
	a AndElem
}
func wrapperAndElem(a AndElem) *andElemWrapper {
	return &andElemWrapper{a}
}
func (c *andElemWrapper) makeCond(cb condBuilder) condBuilder {
    q, v := c.a.mkAndElem()
    if len(q) > 0 {
        cb = cb.appendCond(q, v...)
    }
    return cb
}
func (c *andElemWrapper) mkAndElem() (string, []interface{}) { return "", nil }

type dummyAndElem struct{}
func (a *dummyAndElem) makeCond(cb condBuilder) condBuilder { return cb }
func (a *dummyAndElem) mkAndElem() (string, []interface{}) { return "", nil }

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
type onlyCond struct {
	dummyAndElem
	cond string
}
func (e *onlyCond) mkAndElem() (string, []interface{}) {
	if len(e.cond) == 0 {
		return "", nil
	}
	return e.cond, nil
}

type eqCond struct {
	dummyAndElem
	field string
	val interface{}
}
func (e *eqCond) mkAndElem() (string, []interface{}) {
	if len(e.field) == 0 {
		return "", nil
	}
	backquote := getQuote(e.field)
	return fmt.Sprintf("%s%s%s=?", backquote, e.field, backquote), []interface{}{e.val}
}

type eqExprCond struct {
	dummyAndElem
	field string
	expr string
}
func (e *eqExprCond) mkAndElem() (string, []interface{}) {
	if len(e.field) == 0 || len(e.expr) == 0 {
		return "", nil
	}
	backquote := getQuote(e.field)
	return fmt.Sprintf("%s%s%s=%s", backquote, e.field, backquote, e.expr), nil
}

type opCond struct {
	dummyAndElem
	field string
	op  string
	val interface{}
}
func (e *opCond) mkAndElem() (string, []interface{}) {
	if len(e.field) == 0 {
		return "", nil
	}
	if len(e.op) == 0 {
		// e.op = "="
		// neglect val
		return e.field, nil
	}
	backquote := getQuote(e.field)
	return fmt.Sprintf("%s%s%s %s ?", backquote, e.field, backquote, e.op), []interface{}{e.val}
}

type opExprCond struct {
	dummyAndElem
	field string
	op  string
	expr string
}
func (e *opExprCond) mkAndElem() (string, []interface{}) {
	if len(e.field) == 0 {
		return "", nil
	}
	if len(e.op) == 0 || len(e.expr) == 0 {
		// e.op = "="
		// neglect val
		return e.field, nil
	}
	backquote := getQuote(e.field)
	return fmt.Sprintf("%s%s%s %s %s", backquote, e.field, backquote, e.op, e.expr), nil
}

type andxCond struct {
	dummyAndElem
	conds []AndElem
}
func (e *andxCond) mkAndElem() (string, []interface{}) {
	return joinAndElems(e.conds, "AND")
}

type orxCond struct {
	dummyAndElem
	conds []AndElem
}
func (e *orxCond) mkAndElem() (string, []interface{}) {
	return joinAndElems(e.conds, "OR")
}

type notCond struct {
	dummyAndElem
	conds []AndElem
}
func (e *notCond) mkAndElem() (string, []interface{}) {
	q, v := joinAndElems(e.conds, "NOT")
	if len(q) == 0 {
		return q, v
	}
	return fmt.Sprintf("NOT (%s)", q), v
}

type inCond struct {
	dummyAndElem
	field string
	val []interface{}
}
func (i *inCond) mkAndElem() (string, []interface{}) {
	return makeInElem(i.field, i.val, "IN")
}

type notInCond struct {
	dummyAndElem
	field string
	val []interface{}
}
func (i *notInCond) mkAndElem() (string, []interface{}) {
	return makeInElem(i.field, i.val, "NOT IN")
}

type sqlCond struct {
	sql string
}
func (s *sqlCond) makeCond(cb condBuilder) condBuilder {
	if sess1, ok := cb.(*xormSession); ok {
		return (*xormSession)((*Session)(sess1).Sql(s.sql))
	}
	return cb
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
	field []string
}
func (o *groupBy) makeBy(sess *Session) *Session {
	return sess.GroupBy(strings.Join(o.field, ","))
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
