package dbx

import (
	"github.com/rosbit/xorm"
)

type (
	Session = xorm.Session
	FnIterate = xorm.IterFunc
)

type (
	StmtResult interface {}

	Stmt interface {
		Exec(bean interface{}) (StmtResult, error)
		setSession(session *Session)
	}

	Cond interface {
		makeCond(sess *Session) *Session
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

	Options struct {
		bys []by
		limit limit
		session *Session
	}

	O func(opts *Options)
)

// --- transaction ---
type (
	ArgKey = string
	TxStmt = Pipe
	TxStep = Bolt

	TxA func(args *map[ArgKey]interface{})
)

