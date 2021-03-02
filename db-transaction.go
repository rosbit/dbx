package dbx

// --- transaction ---

import (
	"reflect"
	"fmt"
)

func (f1 TxStepHandler) sameHandler(f2 TxStepHandler) bool {
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}

var (
	// special tx step handlers
	CommitBeforeExecStmt   = func(*TxPrevStepRes)(*TxNextStep, error) { return nil, nil }
	RollbackBeforeExecStmt = func(*TxPrevStepRes)(*TxNextStep, error) { return nil, nil }
	CommitAfterExecStmt TxStepHandler = nil
)

func RunTx(firstStep *TxNextStep, stepHandlers StepHandlers) (err error) {
	db := getDefaultConnection()
	return db.RunTx(firstStep, stepHandlers)
}

// run a transation step by step
func (db *DBI) RunTx(firstStep *TxNextStep, stepHandlers StepHandlers) (err error) {
	if firstStep == nil || stepHandlers == nil || len(stepHandlers) == 0 {
		return fmt.Errorf("bad request for calling RunTx")
	}

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); ok {
				return // err to RunTx() caller
			}
			fmt.Printf("panic in RunTx: %v\n", r)
		}
	}()

	session := db.NewSession()
	defer session.Close()

	if err = session.Begin(); err != nil {
		return err
	}

	nextStep := firstStep

	for {
		stmt, bean, step := nextStep.Stmt, nextStep.Bean, nextStep.Step
		handleTxStep, ok := stepHandlers[step]
		if !ok {
			return fmt.Errorf("no tx handler found for step #%d", step)
		}
		if handleTxStep != nil {
			if handleTxStep.sameHandler(CommitBeforeExecStmt) {
				session.Commit()
				return nil
			} else if handleTxStep.sameHandler(RollbackBeforeExecStmt) {
				return nil
			}
		}

		res, err := stmt.Exec(bean, session)
		if err != nil {
			return err
		}
		if handleTxStep.sameHandler(CommitAfterExecStmt) {
			session.Commit()
			return nil
		}

		nextStep, err = handleTxStep(&TxPrevStepRes{Session:session, Step:step, Bean:bean, Res:res, ExArgs:nextStep.ExArgs})
		if err != nil {
			return err
		}
		if nextStep == nil || nextStep.Stmt == nil {
			session.Commit()
			return nil
		}
	}
}
