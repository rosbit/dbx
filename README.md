# dbx

 an ORM utitlity to handle database.

## Usage
 - Import pacakge
   ```go
   import "github.com/rosbit/dbx"
   ```

 - Create db instance
   ```go
   db, err := dbx.CreateDBInstance(dataSourceName, debug)
   db, err := dbx.CreateDBDriverConnection("mysql", dataSourceName, debug)
   ```

 - Query
   ```go
   type User struct {
       Id int
       Name string
   }

   var user User
   has, err := db.GetOne("user", []dbx.Eq{dbx.NewEq("id", 1)}, nil, &user)

   res, err := db.NewQueryStmt("user", []dbx.Eq{dbx.NewEq("id", 1)}, nil, nil, nil).Exec(&user, nil)

   var users []User
   err := db.Find("user", []dbx.Eq{dbx.NewOp("id", ">", 1)}, nil, 0, 10, &users)

   err := db.Select("user", []string{"id","name"}, &users)

   err := db.SQL("user", "select id,name from user", &users)

   c, err := db.Iter("user", []dbx.Eq{dbx.NewOp("id", ">=", 1)}, nil, &user)
   if err == nil {
       for u := range c {
            fmt.Printf("%v\n", u)
       }
   }
   ```

 - Insert/Update/Delete
   ```go
   user := User{
      Id: 0,
      Name: "hi",
   }
   _, err := db.NewInsertStmt("user").Exec(&user, nil)

   user.Name = "haha"
   _, err := db.NewUpdateStmt("user", []dbx.Eq{dbx.NewEq("id", 1)}, nil, []string{"name"}).Exec(&user, nil)

   _, err := db.NewDeleteStmt("user", []dbx.Eq{dbx.NewEq("id", 1)}, nil).Exec(&user, nil)
   ```

 - Transanction
   ```go
   type Balance {
      UserId int
      Balance int
   }

   // transaction steps
   const (
      tx_find_user = iota
      tx_find_balance
      tx_inc_balance
   )

   func IncUserBalance(userId int, balance int) error {
     firstStep := &dbx.TxNextStep {
     Step: tx_find_user,
       Stmt: db.NewQueryStmt("user", []dbx.Eq{dbx.NewEq("id", userId)}, nil, nil, nil),
       Bean: &User{},
       ExArgs: []interface{}{balance, userId},
     }
  
     // call RunTx to run a transaction. Commit if no error ocurrs, otherwise it will rollback. 
     return db.RunTx(firstStep, map[int]dbx.TxStepHandler {
        tx_find_user: user_found,
        tx_find_balance: balance_found,
        tx_inc_balance: dbx.CommitAfterExecStmt,
     })
   }

   // --- step handler ---
   func user_found(step *TxPrevStepRes) (*TxNextStep, error) {
      has := step.Res.(bool)
      if !has {
         return nil, fmt.Error("user not found")
      }
     
      user := step.Bean.(*User)
      exArgs := step.ExArgs
      return &dbx.TxNextStep{
   	   Step: tx_find_balance,
		   Stmt: db.NewQueryStmt("balance", []dbx.Eq{dbx.NewEq("user_id", user.Id)}, nil, nil, nil),
		   Bean: &Balance{},
		   ExArgs: exArgs,
	   }, nil
   }
   
   func balance_found(step *TxPrevStepRes) (*TxNextStep, error) {
      exArgs := step.ExArgs
      incBalance := exArgs[0].(int)
      userId := exArgs[1].(int)
      has := step.Res.(bool)
      if !has {
          // insert a new one
          return &dbx.TxNextStep{
             Step: tx_inc_balance,
             Stmt: db.NewInsertStmt("balance"),
             Bean: &Balance{UserId: userId, Balance: incBalance},
          }, nil
      }

      // increment balance, update it
      balance := step.Bean.(*Balance)
      balance += incBalance
      return &dbx.TxNextStep{
          Step: tx_inc_balance,
          Stmt: db.NewUpdateStmt("balance", []dbx.Eq{dbx.NewEq("user_id", userId)}, nil, []string{"balance}),
          Bean: balance,
      }, nil
   }
   ```
