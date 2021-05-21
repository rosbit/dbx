package dbx

// --- transaction ---

import (
	"reflect"
	"fmt"
)

var (
	// special tx step handlers
	CommitAfterExecStmt TxStepHandler = nil
)

func (f1 TxStepHandler) sameHandler(f2 TxStepHandler) bool {
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}

// --- args options ----

func TxArg(name ArgKey, val interface{}) TxA {
	return func(args *map[ArgKey]interface{}) {
		if len(name) > 0 {
			(*args)[name] = val
		}
	}
}

func TxArgs(args map[ArgKey]interface{}) TxA {
	return func(oldArgs *map[ArgKey]interface{}) {
		if len(args) == 0 {
			return
		}

		if len(*oldArgs) == 0 {
			*oldArgs = args
			return
		}

		for k, v := range args {
			(*oldArgs)[k] = v
		}
	}
}

func TxCopyArgs(step *TxStepRes) TxA {
	return TxArgs(step.args)
}

// --- create next step with optional args ---
func NextStep(handler TxStepHandler, stmt Stmt, bean interface{}, txArgs ...TxA) *TxStep {
	args := make(map[ArgKey]interface{})
	for _, txArg := range txArgs {
		txArg(&args)
	}

	return &TxStep{handler, stmt, bean, args}
}

// --- previous step ---
func (step *TxStepRes) Has() bool {
	return step.res.(bool)
}

func (step *TxStepRes) Val() interface{} {
	return step.bean
}

func (step *TxStepRes) Arg(key ArgKey) (arg interface{}) {
	if len(step.args) == 0 {
		return nil
	}
	arg, _ = step.args[key]
	return
}

func (step *TxStepRes) Session() (*Session) {
	return step.session
}

func (step *TxStepRes) DB() (*DBI) {
	return step.db
}

// run a transation step by step
func RunTx(firstStep *TxStep) (err error) {
	db := getDefaultConnection()
	return db.RunTx(firstStep)
}

func (db *DBI) RunTx(firstStep *TxStep) (err error) {
	if firstStep == nil || firstStep.stmt == nil {
		return nil
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
		stmt, bean, handleTxStep, args := nextStep.stmt, nextStep.bean, nextStep.step, nextStep.args
		stmt.setSession(session)
		res, err := stmt.Exec(bean)
		if err != nil {
			return err
		}

		if handleTxStep.sameHandler(CommitAfterExecStmt) {
			session.Commit()
			return nil
		}

		if nextStep, err = handleTxStep(&TxStepRes{session, db, bean, res, args}); err != nil {
			return err
		}

		if nextStep == nil || nextStep.stmt == nil {
			session.Commit()
			return nil
		}
	}
}
