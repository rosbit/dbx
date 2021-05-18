package dbx

import (
	"reflect"
)

func Where(eq ...Cond) []Cond {
	return eq
}

func Cols(field ...string) []string {
	return field
}

func Eq(fieldName string, val interface{}) AndElem {
	return &eqCond{fieldName, val}
}

// op: "=", "!=", "<>", ">", ">=", "<", "<=", "like"
func Op(fieldName string, op string, val interface{}) AndElem {
	return &opCond{fieldName, op, val}
}

func And(cond ...AndElem) AndElem {
	return &andxCond{cond}
}

func Or(cond ...AndElem) AndElem {
	return &orxCond{cond}
}

func Not(cond ...AndElem) AndElem {
	return &notCond{cond}
}

func In(fieldName string, val ...interface{}) AndElem {
	return &inCond{fieldName, copy_i(val...)}
}

func NotIn(fieldName string, val ...interface{}) AndElem {
	return &notInCond{fieldName, copy_i(val...)}
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
