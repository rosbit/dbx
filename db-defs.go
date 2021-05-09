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
		Exec(bean interface{}, session ...*Session) (StmtResult, error)
	}

	Cond interface {
		makeCond(sess *Session) *Session
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
	}

	O func(opts *Options)
)

// --- transaction ---
type (
	StepKey = string
	ArgKey  = string

	TxStepRes struct {
		session *Session
		db      *DBI
		step     StepKey
		bean     interface{}
		res      StmtResult
		args     map[ArgKey]interface{}
	}

	TxStep struct {
		step  StepKey
		stmt  Stmt
		bean  interface{}
		args  map[ArgKey]interface{}
	}

	TxStepHandler func(*TxStepRes)(*TxStep, error)
	StepHandlers map[StepKey]TxStepHandler

	TxA func(args *map[ArgKey]interface{})
)

