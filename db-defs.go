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
		Exec(bean interface{}, session *Session) (StmtResult, error)
	}

	Eq interface {
		makeCond(sess *Session) *Session
	}
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
