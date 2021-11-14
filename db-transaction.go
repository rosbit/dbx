package dbx

import (
	"fmt"
)

type TxStmt struct {
	*dbxStmt
	args map[ArgKey]interface{}
	session *Session
}

func TxStmts(stmts ...FnTxStmt) []FnTxStmt {
	return stmts
}

func (db *DBI) newTxStmt(session *Session, table string, txArgs ...TxA) *TxStmt {
	args := make(map[ArgKey]interface{})
	for _, txArg := range txArgs {
		txArg(&args)
	}

	return &TxStmt{
		dbxStmt: db.XStmt(table).XSession(session),
		args: args,
		session: session,
	}
}

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

func (ts *TxStmt) Arg(key ArgKey) (arg interface{}) {
	if len(ts.args) == 0 {
		return nil
	}
	arg, _ = ts.args[key]
	return
}

func (ts *TxStmt) Set(name ArgKey, val interface{}) {
	if len(name) == 0 {
		return
	}
	if ts.args == nil {
		ts.args = map[ArgKey]interface{}{name: val}
		return
	}
	ts.args[name] = val
}

func Tx(stmts []FnTxStmt, txArgs ...TxA) error {
	db := getDefaultConnection()
	return db.Tx(stmts, txArgs...)
}

func (db *DBI) Tx(stmts []FnTxStmt, txArgs ...TxA) (err error) {
	if len(stmts) == 0 {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); ok {
				return // err to call FnTxStmt()
			}
			fmt.Printf("panic in FnTxStmt: %v\n", r)
		}
	}()

	session := db.NewSession()
	defer session.Close()
	if err = session.Begin(); err != nil {
		return err
	}

	txStmt := db.newTxStmt(session, "", txArgs...)
	for i, _ := range stmts {
		fnTx := stmts[i]
		if fnTx == nil {
			continue
		}
		if err = fnTx(txStmt); err != nil {
			return err
		}
	}

	session.Commit()
	return nil
}
