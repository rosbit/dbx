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
  dataSourceName := dbx.GenerateMysqlDSN(dbx.Host("127.0.0.1"), dbx.User("root"), dbx.DBName("test"))
  dataSourceName := dbx.GenerateMysqlDSN(dbx.Host("127.0.0.1", 3306), dbx.User("root", ""))
  dataSourceName := dbx.GenerateMysqlDSN(dbx.DomainSocket("/tmp/mysql.sock"))
  
  err := dbx.CreateMysqlConnection(dataSourceName, debug)
  db, err := dbx.CreateMysqlInstance(dataSourceName, debug)
  db, err := dbx.CreateDBDriverConnection("mysql", dataSourceName, debug)
  ```

- Statements
  
  ```go
  type User struct {
      Id int
      Name string
      Age int
  }
  
  var user User
  var users []User
  
  // statement
  //  Query
  has, err := db.XStmt("user").Where(dbx.Eq("name", "rosbit")).Get(&user)
  err := db.XStmt("user").Or(dbx.Eq("name", "rosbit"), dbx.Eq("age", 1)).Desc("name").Get(&user)
  err := db.XStmt("user").Or(dbx.Eq("name", "rosbit"), dbx.Eq("age", 1)).Limit(2).List(&users)
  
  //  iterate
  for uu := range db.XStmt("user").Or(dbx.Eq("name", "rosbit"), dbx.Eq("age", 1)).Iter(&user) {
      u := uu.(*User)
      // do something with u
  }
  
  //  insert/update/delete
  err := db.XStmt("user").Insert(&user)
  affected, err := db.XStmt("user").Where(dbx.Eq("id", user.Id)).Cols("name", "age").Update(&user)
  err := db.XStmt("user").Where(dbx.Eq("id", user.Id)).Delete(&user)
  
  count, err := db.XStmt("user").Where(dbx.Eq("name", "rosbit")).Count(&user)
  sum, err := db.XStmt("user").Where(dbx.Eq("name", "rosbit")).Sum(&user, "age")
  ```

- Join
  
  ```go
  type User struct {
      Id int
      Name string
      Age int
  }
  type Tag struct {
      Id int
      Name string
  }
  type UserTag struct {
      Id int
      UserId int
      TagId int
  }
  
  type UserWithTags struct {
      User `xorm:"extends"`
      UserTag `xorm:"extends"`
      Tag `xorm:"extends"`
  }
  
  var userWithTags []UserWithTags
  err := dbx.XStmt().InnerJoin("user", "userTag", "user.id=userTag.user_id").
                     NextInnerJoin("tag", "userTag.tag_id=tag.id").
             Where(dbx.Eq("user.id", 1)).
             List(&userWithTags)
  ```

- Transaction
  
  ```go
  type Balance struct {
     UserId int
     Balance int
  }
  
  const (
     // args name
     arg_balance = "balance"
     arg_user_id = "user_id"
  )
  
  func IncUserBalance(db *dbx.DBI, userId int, balance int) error {
    // call Tx to run a transaction. Commit if no error ocurrs, otherwise it will rollback. 
    return db.Tx(
       dbx.TxStmts(
           find_user,
           inc_balance,
       ),
       dbx.TxArg(arg_balance, balance),
       dbx.TxArg(arg_user_id, userId),
    )
  }
  
  // --- stmt handler ---
  func find_user(stmt *dbx.TxStmt) (error) {
     userId := stmt.Arg(arg_user_id).(int)
     var user User
     has, err := stmt.Table("user").Where(dbx.Eq("id", userId)).Get(&user)
     if err != nil {
        return err
     }
     if !has {
        return fmt.Errorf("user not found")
     }
     return nil
  }
  
  func inc_balance(stmt *dbx.TxStmt) (error) {
     userId := stmt.Arg(arg_user_id).(int)
     incBalance := stmt.Arg(arg_balance).(int)
     var balance Balance
     has, err := stmt.Table("balance").Where(dbx.Eq("user_id", userId)).Get(&balance)
     if err != nil {
        return err
     }
     if !has {
         // insert a new one
         balance.UserId = userId
         balance.Balance = incBalance
  
         return stmt.Table("balance").Insert(&balance)
     }
  
     // increment balance, update it
     balance.Balance += incBalance
     _, err = stmt.Table("balance").Where(dbx.Eq("user_id", userId)).Cols("balance").Update(&balance)
     return err
  }
  ```

- Conditions
  
  ```go
  // conditions can be grouped by dbx.Where() or XStmt::Where()
  
  // And
  dbx.Eq("a", 1)
  dbx.Op("b", ">", 2))  // -> where
  
  dbx.Where(dbx.Eq("a", 1), dbx.Op("b", ">", 2))
  dbx.Where(dbx.And(dbx.Eq("a", 1), dbx.Op("b", ">", 2)))
  
  // Or
  dbx.Or(
     dbx.Eq("a", 1),
     dbx.Op("b", "<", 2),
     dbx.Op("c", ">=", 3),
  ) // -> where
  dbx.Where(
     dbx.Or(
        dbx.Eq("a", 1),
        dbx.Op("b", "<", 2),
        dbx.Op("c", ">=", 3),
     ),
  )
  
  // NOT
  //  NOT AND
  dbx.Not(
     dbx.Eq("a", 1),
     dbx.Op("b", ">", 2),
  )
  dbx.Not(
     dbx.And(
       dbx.Eq("a", 1),
       dbx.Op("b", ">", 2),
     ),
  ) // -> where
  dbx.Where(dbx.Not(dbx.Eq("a", 1), dbx.Op("b", ">", 2)))
  dbx.Where(dbx.Not(dbx.And(dbx.Eq("a", 1), dbx.Op("b", ">", 2))))
  // NOT OR
  dbx.Not(
     dbx.Or(
       dbx.Eq("a", 1),
       dbx.Op("b", "<", 2),
       dbx.Op("c", ">=", 3),
     ),
  ) // -> where
  dbx.Where(dbx.Not(dbx.Or(dbx.Eq("a", 1), dbx.Op("b", "<", 2), dbx.Op("c", ">=", 3))))
  
  // IN
  dbx.In("id", 1, 3, 5)
  dbx.In("id", []int{1, 3, 5}) // -> where
  
  dbx.Where(dbx.In("id", 1, 3, 5))
  dbx.Where(dbx.In("id", []int{1, 3, 5}))
  
  // not IN
  dbx.NotIn("id", 1, 3, 5)
  dbx.Not(dbx.In("id", []int{1, 3, 5}))
  dbx.NotIn("id", 1, 3, 5)
  dbx.Not(dbx.In("id", []int{1, 3, 5})) // -> where
  
  dbx.Where(dbx.NotIn("id", 1, 3, 5))
  dbx.Where(dbx.Not(dbx.In("id", []int{1, 3, 5})))
  dbx.Where(dbx.NotIn("id", 1, 3, 5))
  dbx.Where(dbx.Not(dbx.In("id", []int{1, 3, 5})))
  
  // SQL
  dbx.Where(dbx.Sql("select id,name from user"))
  ```

## Status

The package is fully tested.

## Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.

__Convention:__ fork the repository and make changes on your fork in a feature branch.
