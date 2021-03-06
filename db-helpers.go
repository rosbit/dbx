package dbx

import (
	"fmt"
)

func Conds(eq ...Cond) []Cond {
	return eq
}

func Bys(by ...By) []By {
	return by
}

func Cols(field ...string) []string {
	return field
}

func NewEq(fieldName string, val ...interface{}) Cond {
	return &andCond{fieldName, copy_i(val)}
}

func NewConds(cond ...string) Cond {
	return &andxCond{cond}
}

func OrderByDesc(field ...string) By {
	return &descOrderBy{field}
}

func OrderByAsc(field ...string) By {
	return &ascOrderBy{field}
}

func GroupBy(field string) By {
	return &groupBy{field}
}

func LimitCount(count int, offset ...int) Limit {
	if count <= 0 {
		return nil
	}
	l := &limitT{count:count}
	if len(offset) > 0 && offset[0] >= 0 {
		l.offset = offset[0]
	}
	return l
}

// op: "=", "!=", "<>", ">", ">=", "<", "<=", "like"
func NewOp(fieldName string, op string, val ...interface{}) Cond {
	return &opCond{fieldName, op, copy_i(val)}
}

// f1, v1, f2, v2, ... => (f1=v1 OR f2=v2 OR ...)
func NewOr(fieldName string, val ...interface{}) Cond {
	fields := []string{fieldName}
	vals := []interface{}{}

	for i, v := range val {
		if i % 2 == 0 {
			vals = append(vals, v)
		} else {
			fields = append(fields, fmt.Sprintf("%s", v))
		}
	}
	return &orCond{fields, vals}
}

func NewIn(fieldName string, val ...interface{}) Cond {
	return &inCond{fieldName, copy_i(val)}
}

func NewNotIn(fieldName string, val ...interface{}) Cond {
	return &notInCond{fieldName, copy_i(val)}
}

func copy_i(vals ...interface{}) []interface{} {
	if len(vals) == 0 {
		return nil
	}
	vs := make([]interface{}, len(vals))
	for i, _ := range vals {
		vs[i] = vals[i]
	}
	return vs
}

func NewSql(sql string) Cond {
	return &sqlCond{sql}
}

