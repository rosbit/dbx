package dbx

import (
	"github.com/rosbit/xorm"
)

type (
	Session = xorm.Session
	FnIterate = xorm.IterFunc
	xormSession xorm.Session // adding methods is allowed
)

type (
	StmtResult interface {}

	Stmt interface {
		Exec(bean interface{}) (StmtResult, error)
		setSession(session *Session)
	}

	condBuilder interface {
		appendCond(query string, args ...interface{}) condBuilder
	}

	Cond interface {
		makeCond(condBuilder) condBuilder
	}

	AndElem interface {
		Cond
		mkAndElem() (q string, v []interface{})
	}

	by interface {
		makeBy(sess *Session) *Session
	}

	limit interface {
		makeLimit(sess *Session) *Session
	}

	Set interface {
		makeSetClause() (setClaus string, v interface{})
	}

	Options struct {
		bys []by
		limit limit
		session *Session
		selection string
	}

	O func(opts *Options)
)

// --- transaction ---
type (
	ArgKey = string
	TxA func(args *map[ArgKey]interface{})
	FnTxStmt func(*TxStmt) error
)

