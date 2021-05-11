# dbx

 an ORM utitlity to handle database.

## Usage
 - Import pacakge
   ```go
   import "github.com/rosbit/dbx"
   ```

 - Create db instance
   ```go
   dataSourceName := dbx.GenerateMysqlDSN(dbx.DBName("test"), dbx.Attr("charset", "utf8mb4"))
   dataSourceName := dbx.GenerateMysqlDSN(dbx.Host("127.0.0.1"), dbx.User("root"))
   dataSourceName := dbx.GenerateMysqlDSN(dbx.Host("127.0.0.1", 3306), dbx.User("root", ""))
   dataSourceName := dbx.GenerateMysqlDSN(dbx.DomainSocket("/tmp/mysql.sock"))
   
   err := dbx.CreateMysqlConnection(dataSourceName, debug)
   db, err := dbx.CreateMysqlInstance(dataSourceName, debug)
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
   has, err := db.GetOne("user", dbx.Where(dbx.Eq("id", 1)), &user)
   has, err := db.GetBy("user", "id", 1, &user)
   res, err := db.QueryStmt("user", dbx.Where(dbx.Eq("id", 1))).Exec(&user)
   
   var users []User
   err := db.List("user", dbx.Where(dbx.Op("id", ">", 1)), &users, dbx.OrderBy("id"), dbx.Limit(10))
   err := db.List("user", dbx.Where(dbx.Op("id", ">", 1)), &users, dbx.OrderByDesc("id"), dbx.Limit(10))
   err := db.Find("user", dbx.Where(dbx.Op("id", ">", 1)), &users, dbx.OrderByDesc("id"), dbx.Limit(10))
   err := db.Select("user", dbx.Cols("id","name"), dbx.Where(dbx.Eq("id", 1)), &users)
   err := db.RunSQL("user", "select id,name from user", &users)
   
   // iterate
   c, err := db.Iter("user", dbx.Where(dbx.Op("id", ">=", 1)), &user)
   if err == nil {
       for u := range c {
            user := u.(*User)
            fmt.Printf("%v\n", user)
       }
   }
   
   if err := db.Iterate("user", dbx.Where(dbx.Op("id", ">=", 1)), &user, func(idx int, bean interface{}){
       fmt.Printf("%v\n", bean.(*User))
   })
   
   // inner join
   type Detail struct {
       Id int
       Detail string
   }
   type UserDetail struct {
       User   `xorm:"extends"`
       Detail `xorm:"extends"`
   }
   var userDetails []UserDetail
   
   if err := db.InnerJoin("user", "detail", "user.id=detail.id", dbx.Where(dbx.Op("user.id", ">", 1)), &userDetails, dbx.Limit(10)); err != nil {
      // xxx
   }
   if _, err := db.InnerJoinStmt("user", "detail", "user.id=detail.id", dbx.Where(dbx.Op("user.id", ">", 1)), dbx.Limit(10)).Exec(&userDetails); err != nil {
      // xxx
   }
   ```
   
 - Insert/Update/Delete
   ```go
   user := User{
      Id: 0,
      Name: "hi",
   }
   err := db.Insert("user", &user)
   _, err := db.InsertStmt("user").Exec(&user)
   
   user.Name = "haha"
   err := db.Update("user", dbx.Where(dbx.Eq("id", 1)), dbx.Cols("name"), &user)
   _, err := db.UpdateStmt("user", dbx.Where(dbx.Eq("id", 1)), dbx.Cols("name")).Exec(&user)
   
   err := db.Delete("user", dbx.Where(dbx.Eq("id", 1)), &user)
   _, err := db.DeleteStmt("user", dbx.Where(dbx.Eq("id", 1))).Exec(&user)
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
   dbx.OrderBy("id", "name")  // equals to
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
   
   const (
      // args name
      arg_balance = "balance"
      arg_user_id = "user_id"
   )
   
   func IncUserBalance(db *dbx.DBI, userId int, balance int) error {
     firstStep := dbx.NextStep(
        user_found,
        db.QueryStmt("user", dbx.Where(dbx.Eq("id", userId))),
        &User{},
        dbx.TxArg(arg_balance, balance),
        dbx.TxArg(arg_user_id, userId),
     )
     
     // call RunTx to run a transaction. Commit if no error ocurrs, otherwise it will rollback. 
     return db.RunTx(firstStep)
   }
   
   // --- step handler ---
   func user_found(step *TxStepRes) (*TxStep, error) {
      if !step.Has() {
         return nil, fmt.Error("user not found")
      }
     
      user := step.Val().(*User)
      return dbx.NextStep(
           balance_found,
           step.DB().QueryStmt("balance", dbx.Where(dbx.Eq("user_id", user.Id))),
           &Balance{},
           dbx.TxCopyArgs(step),
      ), nil
   }
   
   func balance_found(step *TxStepRes) (*TxStep, error) {
      incBalance := step.Arg(arg_balance).(int)
      userId := step.Arg(arg_user_id).(int)
      if !step.Has() {
          // insert a new one
          return dbx.NextStep(
             dbx.CommitAfterExecStmt,
             step.DB().InsertStmt("balance"),
             &Balance{UserId: userId, Balance: incBalance},
          ), nil
      }
   
      // increment balance, update it
      balance := step.Val().(*Balance)
      balance.Balance += incBalance
      return dbx.NextStep(
          dbx.CommitAfterExecStmt,
          step.DB().UpdateStmt("balance", dbx.Where(dbx.Eq("user_id", userId)), dbx.Cols("balance")),
          balance,
      ), nil
   }
   ```
