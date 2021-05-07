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
   has, err := db.Get("user", dbx.Where(dbx.Eq("id", 1)), &user)
   has, err := db.GetBy("user", "id", 1, &user)
   res, err := db.NewQueryStmt("user", dbx.Where(dbx.Eq("id", 1))).Exec(&user)

   var users []User
   err := db.Find("user", dbx.Where(dbx.Op("id", ">", 1)), &users, dbx.OrderByDesc("id"), dbx.Limit(10))
   err := db.Select("user", dbx.Cols("id","name"), dbx.Where(dbx.Eq("id", 1)), &users)
   err := db.RunSQL("user", "select id,name from user", &users)

   c, err := db.Iter("user", dbx.Where(dbx.Op("id", ">=", 1)), &user)
   if err == nil {
       for u := range c {
            user := u.(*User)
            fmt.Printf("%v\n", user)
       }
   }
   ```

 - Insert/Update/Delete
   ```go
   user := User{
      Id: 0,
      Name: "hi",
   }
   err := db.Insert("user", &user)
   _, err := db.NewInsertStmt("user").Exec(&user)

   user.Name = "haha"
   err := db.Update("user", dbx.Where(dbx.Eq("id", 1)), dbx.Cols("name"), &user)
   _, err := db.NewUpdateStmt("user", dbx.Where(dbx.Eq("id", 1)), dbx.Cols("name")).Exec(&user)

   err := db.Delete("user", dbx.Where(dbx.Eq("id", 1)), &user)
   _, err := db.NewDeleteStmt("user", dbx.Where(dbx.Eq("id", 1))).Exec(&user)
   ```

 - Conditions
   ```go
   // conditions can be grouped by dbx.Where()

   // And
   dbx.Eq("a", 1)
   dbx.Op("b", ">", 2))  // -> where

   dbx.Where(dbx.Eq("a", 1), dbx.Op("b", ">", 2))
   dbx.Where(dbx.And("a=1", "b>2"))

   // Or
   dbx.Or("a=1", "b<2", "c>=3") // -> where
   dbx.Where(dbx.Or("a=1", "b<2", "c>=3"))

   dbx.Where(dbx.OrEq("a", 1, "b", 2, "c", 3)) // -> "a=1" "b=2 "c=3"

   // IN
   dbx.In("id", 1, 3, 5)
   dbx.In("id", []int{1, 3, 5}) // -> where

   dbx.Where(dbx.In("id", 1, 3, 5))
   dbx.Where(dbx.In("id", []int{1, 3, 5}))

   // not IN
   dbx.NotIn("id", 1, 3, 5)
   dbx.NotIn("id", []int{1, 3, 5}) // -> where

   dbx.Where(dbx.NotIn("id", 1, 3, 5))
   dbx.Where(dbx.NotIn("id", []int{1, 3, 5}))

   // SQL
   dbx.Where(dbx.Sql("select id,name from user"))
   ```

 - Options
   ```go
   // sorting
   dbx.OrderByDesc("id", "name")
   dbx.OrderByAsc("id", "name")

   // grouping
   dbx.GroupBy("id")

   // limit count
   dbx.Limit(10)
   dbx.Limit(20, 100)  // offset: 100, count: 20
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
       Stmt: db.NewQueryStmt("user", dbx.Where(dbx.Eq("id", userId))),
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
		   Stmt: db.NewQueryStmt("balance", dbx.Where(dbx.Eq("user_id", user.Id))),
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
          Stmt: db.NewUpdateStmt("balance", dbx.Where(dbx.Eq("user_id", userId)), dbx.Cols("balance)),
          Bean: balance,
      }, nil
   }
   ```
