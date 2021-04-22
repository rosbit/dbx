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

	By interface {
		makeBy(sess *Session) *Session
	}

	Limit interface {
		makeLimit(sess *Session) *Session
	}

	Options struct {
		Bys []By
		Count Limit
	}

	O func(opts *Options)
)

// --- transaction ---
type (
	TxPrevStepRes struct {
		Session *Session
		Step     int
		Bean     interface{}
		Res      StmtResult
		ExArgs []interface{}
	}

	TxNextStep struct {
		Step     int
		Stmt
		Bean     interface{}
		ExArgs []interface{}
	}

	TxStepHandler func(*TxPrevStepRes)(*TxNextStep, error)
	StepHandlers map[int]TxStepHandler
)
