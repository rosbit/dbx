package dbx

import (
	"reflect"
	"strings"
)

// -- conditions ---
func Where(eq ...Cond) []Cond {
	return eq
}

func Cols(field ...string) []string {
	return field
}

func OnlyCond(cond string) AndElem {
	return wrapperAndElem(&onlyCond{cond:cond})
}

func Eq(fieldName string, val interface{}) AndElem {
	return wrapperAndElem(&eqCond{field:fieldName, val:val})
}

func Ne(fieldName string, val interface{}) AndElem {
	return Op(fieldName, "<>", val)
}

func Gt(fieldName string, val interface{}) AndElem {
	return Op(fieldName, ">", val)
}

func Ge(fieldName string, val interface{}) AndElem {
	return Op(fieldName, ">=", val)
}

func Lt(fieldName string, val interface{}) AndElem {
	return Op(fieldName, "<", val)
}

func Le(fieldName string, val interface{}) AndElem {
	return Op(fieldName, "<=", val)
}

func EqX(fieldName string, expr string) AndElem {
	return wrapperAndElem(&eqExprCond{field:fieldName, expr:expr})
}

// op: "=", "!=", "<>", ">", ">=", "<", "<=", "like"
func Op(fieldName string, op string, val interface{}) AndElem {
	switch strings.ToLower(op) {
	case "in":
		return In(fieldName, val)
	case "not in":
		return NotIn(fieldName, val)
	default:
		return wrapperAndElem(&opCond{field:fieldName, op:op, val:val})
	}
}

func OpX(fieldName string, op string, expr string) AndElem {
	return wrapperAndElem(&opExprCond{field:fieldName, op:op, expr:expr})
}

func And(cond ...AndElem) AndElem {
	return wrapperAndElem(&andxCond{conds:cond})
}

func Or(cond ...AndElem) AndElem {
	return wrapperAndElem(&orxCond{conds:cond})
}

func Not(cond ...AndElem) AndElem {
	return wrapperAndElem(&notCond{conds:cond})
}

func In(fieldName string, val ...interface{}) AndElem {
	return wrapperAndElem(&inCond{field:fieldName, val:copy_i(val...)})
}

func NotIn(fieldName string, val ...interface{}) AndElem {
	return wrapperAndElem(&notInCond{field:fieldName, val:copy_i(val...)})
}

func Sql(sql string) Cond {
	return &sqlCond{sql}
}

func copy_i(vals ...interface{}) []interface{} {
	switch len(vals) {
	case 0:
		// no args
		return nil
	case 1:
		// only 1 arg
		if vals[0] == nil {
			// the arg is nil
			return nil
		}

		// whether the arg is array or slice
		val := reflect.ValueOf(vals[0])
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
		default:
			// not array/slice, wrap it as a slice
			return []interface{}{vals[0]}
		}

		// convert array/slice into []interface{}
		arrLen := val.Len()
		if arrLen == 0 {
			return nil
		}
		vs := make([]interface{}, arrLen)
		for i:=0; i<arrLen; i++ {
			vs[i] = val.Index(i).Interface()
		}
		return vs
	default:
		// args is []interface{}, return it directly
		return vals
	}
}

// --- order by, group by, limit ---
var OrderBy = OrderByDesc

func OrderByDesc(field ...string) O {
	return func(opts *Options) {
		opts.bys = append(opts.bys, &descOrderBy{field})
	}
}

func OrderByAsc(field ...string) O {
	return func(opts *Options) {
		opts.bys = append(opts.bys, &ascOrderBy{field})
	}
}

func GroupBy(field ...string) O {
	return func(opts *Options) {
		opts.bys = append(opts.bys, &groupBy{field})
	}
}

func Limit(count int, offset ...int) O {
	return func(opts *Options) {
		if count <= 0 {
			return
		}

		l := &limitOffset{count:count}
		if len(offset) > 0 && offset[0] >= 0 {
			l.offset = offset[0]
		}
		opts.limit = l
	}
}

func WithSession(session *Session) O {
	return func(opts *Options) {
		opts.session = session
	}
}

func SelectCols(selection string) O {
	return func(opts *Options) {
		opts.selection = selection
	}
}

// --- update SET claus --
func Sets(sets ...Set) []Set {
	return sets
}

func SetValue(field string, val interface{}) Set {
	return &setValue{field, val}
}

func SetExpr(field string, expr string) Set {
	return &setExpr{field, expr}
}
