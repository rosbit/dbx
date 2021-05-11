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

func TxCopyArgs(step *TxStepRes) TxA {
	return func(args *map[ArgKey]interface{}) {
		if len(*args) == 0 {
			*args = step.args
			return
		}

		for k, v := range step.args {
			(*args)[k] = v
		}
	}
}

// --- create next step with optional args ---
func NextStep(name StepKey, stmt Stmt, bean interface{}, txArgs ...TxA) *TxStep {
	args := make(map[ArgKey]interface{})
	for _, txArg := range txArgs {
		txArg(&args)
	}

	return &TxStep{name, stmt, bean, args}
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
func RunTx(firstStep *TxStep, stepHandlers StepHandlers) (err error) {
	db := getDefaultConnection()
	return db.RunTx(firstStep, stepHandlers)
}

func (db *DBI) RunTx(firstStep *TxStep, stepHandlers StepHandlers) (err error) {
	if firstStep == nil || len(stepHandlers) == 0 {
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
		if nextStep == nil || nextStep.stmt == nil {
			session.Commit()
			return nil
		}

		stmt, bean, step, args := nextStep.stmt, nextStep.bean, nextStep.step, nextStep.args
		handleTxStep, ok := stepHandlers[step]
		if !ok {
			return fmt.Errorf("no tx handler found for step %s", step)
		}

		res, err := stmt.Exec(bean, session)
		if err != nil {
			return err
		}

		if handleTxStep.sameHandler(CommitAfterExecStmt) {
			session.Commit()
			return nil
		}

		if nextStep, err = handleTxStep(&TxStepRes{session, db, step, bean, res, args}); err != nil {
			return err
		}
	}
}
