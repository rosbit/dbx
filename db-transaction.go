package dbx

import (
	"fmt"
)

type Pipe struct {
	*dbxStmt
	args map[ArgKey]interface{}
	session *Session
}

type FnBolt func(*Pipe)(*Bolt, error)

type Bolt struct {
	pipe *Pipe
	bolt FnBolt
}

func (ts *Pipe) Next(txArgs ...TxA) *Pipe {
	return ts.engine.newPipe(ts.session, ts.table, txArgs...)
}

func NextBolt(stmt *Pipe, bolt FnBolt) *Bolt {
	return &Bolt{
		pipe: stmt,
		bolt: bolt,
	}
}

var (
	NextStep = NextBolt
	RunTx = PipeTx
)

func (db *DBI) newPipe(session *Session, table string, txArgs ...TxA) *Pipe {
	args := make(map[ArgKey]interface{})
	for _, txArg := range txArgs {
		txArg(&args)
	}

	return &Pipe{
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

func (ts *Pipe) CopyArgs() TxA {
	return TxArgs(ts.args)
}

func (ts *Pipe) Arg(key ArgKey) (arg interface{}) {
	if len(ts.args) == 0 {
		return nil
	}
	arg, _ = ts.args[key]
	return
}

func PipeTx(bolt FnBolt, txArgs ...TxA) (err error) {
	db := getDefaultConnection()
	return db.PipeTx(bolt, txArgs...)
}

func (db *DBI) RunTx(bolt FnBolt, txArgs ...TxA) (err error) {
	return db.PipeTx(bolt, txArgs...)
}

func (db *DBI) PipeTx(bolt FnBolt, txArgs ...TxA) (err error) {
	if bolt == nil {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); ok {
				return // err to BoltTx() caller
			}
			fmt.Printf("panic in BoltTx: %v\n", r)
		}
	}()

	session := db.NewSession()
	defer session.Close()
	if err = session.Begin(); err != nil {
		return err
	}

	handleBolt := bolt
	stmt := db.newPipe(session, "", txArgs...)
	for {
		nextBolt, err := handleBolt(stmt)
		if err != nil {
			return err
		}
		if nextBolt == nil {
			session.Commit()
			return nil
		}
		stmt, handleBolt = nextBolt.pipe, nextBolt.bolt
		if stmt == nil || handleBolt == nil {
			session.Commit()
			return nil
		}
	}
}
