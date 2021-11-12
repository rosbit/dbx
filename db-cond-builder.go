// implementation of interface condBuilder
package dbx

import (
	"strings"
)

func buildConds(cb condBuilder, conds []Cond) {
	for i, _ := range conds {
		conds[i].makeCond(cb)
	}
}

// xormSession::appendCond
func (sess *xormSession) appendCond(q string, v ...interface{}) condBuilder {
	sess1 := (*Session)(sess).And(q, v...)
	return (*xormSession)(sess1)
}

// sqlBuilder, only used for updateSetStmt
type sqlBuilder struct {
	q *strings.Builder
	v []interface{}
	hasWhere bool
}

func newSqlBuilder() *sqlBuilder {
	return &sqlBuilder{
		q: &strings.Builder{},
		v: []interface{}{""}, // v[0] is the place holder for SQL
	}
}

// sqlBuilder::appendCond
func (sb *sqlBuilder) appendCond(q string, v ...interface{}) condBuilder {
	if len(q) == 0 {
		return sb
	}

	if sb.hasWhere {
		sb.q.WriteString(" AND ")
	} else {
		sb.hasWhere = true
		sb.q.WriteString(" ")
	}
	sb.q.WriteString("(")
	sb.q.WriteString(q)
	sb.q.WriteString(")")
	if len(sb.v) == 0 {
		sb.v = v
	} else {
		sb.v = append(sb.v, v...)
	}
	return sb
}

func (sb *sqlBuilder) appendSets(sets []Set) {
	firstClause := true

	for i, _ := range sets {
		c, v := sets[i].makeSetClause()
		if len(c) == 0 {
			continue
		}

		if firstClause {
			sb.q.WriteString(" ")
			firstClause = false
		} else {
			sb.q.WriteString(",")
		}
		sb.q.WriteString(c)
		if v != nil {
			sb.v = append(sb.v, v)
		}
	}
}

func (sb *sqlBuilder) toParams() []interface{} {
	sb.v[0] = sb.q.String() // replace the SQL place holder
	return sb.v
}

