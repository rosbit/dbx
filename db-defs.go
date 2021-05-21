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
	ArgKey  = string

	TxStepHandler func(*TxStepRes)(*TxStep, error)

	TxStepRes struct {
		session *Session
		db      *DBI
		bean     interface{}
		res      StmtResult
		args     map[ArgKey]interface{}
	}

	TxStep struct {
		step  TxStepHandler
		stmt  Stmt
		bean  interface{}
		args  map[ArgKey]interface{}
	}

	TxA func(args *map[ArgKey]interface{})
)

