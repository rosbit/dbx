package dbx

import (
	"fmt"
)

// implementation of interface Set
type setValue struct {
	field string
	val interface{}
}
func (s *setValue) makeSetClause() (setClaus string, v interface{}) {
	if len(s.field) == 0 {
		return
	}
	backquote := getQuote(s.field)
	setClaus = fmt.Sprintf("%s%s%s=?", backquote, s.field, backquote)
	v = s.val
	return
}

type setExpr struct {
	field string
	expr string
}
func (s *setExpr) makeSetClause() (setClaus string, v interface{}) {
	if len(s.field) == 0 {
		return
	}
	backquote := getQuote(s.field)
	setClaus = fmt.Sprintf("%s%s%s=%s", backquote, s.field, backquote, s.expr)
	return
}

