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

func Eq(fieldName string, val ...interface{}) Cond {
	return &andCond{fieldName, copy_i(val)}
}

func And(cond ...string) Cond {
	return &andxCond{cond}
}

// op: "=", "!=", "<>", ">", ">=", "<", "<=", "like"
func Op(fieldName string, op string, val ...interface{}) Cond {
	return &opCond{fieldName, op, copy_i(val)}
}

// f1, v1, f2, v2, ... => (f1=v1 OR f2=v2 OR ...)
func Or(fieldName string, val ...interface{}) Cond {
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

func In(fieldName string, val ...interface{}) Cond {
	return &inCond{fieldName, copy_i(val)}
}

func NotIn(fieldName string, val ...interface{}) Cond {
	return &notInCond{fieldName, copy_i(val)}
}

func Sql(sql string) Cond {
	return &sqlCond{sql}
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

