package dbx

import (
	_ "github.com/go-sql-driver/mysql"
	"bytes"
	"fmt"
)

// create a default instance of mysql connection
func CreateMysqlConnection(dsn string, debug bool) (err error) {
	DB, err = CreateMysqlInstance(dsn, debug)
	return
}

// create an instance of mysql connection with an dsn
func CreateMysqlInstance(dsn string, debug bool) (db *DBI, err error) {
	return CreateDriverDBInstance("mysql", dsn, debug)
}

// generate mysql DSN
func GenerateMysqlDSN(options ...mysqlOption) string {
	params := &mysqldParams{
		protocol: "tcp",
		host: "localhost:3306",
		user: "root",
		attrs: map[string]interface{}{
			"charset": "utf8mb4",
		},
	}

	for _, o := range options {
		o(params)
	}

	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s", params.user, params.passwd, params.protocol, params.host, params.db)
	if len(params.attrs) > 0 {
		first := true
		p := bytes.NewBufferString(dsn)
		for k, v := range params.attrs {
			if first {
				fmt.Fprintf(p, "?%s=%v", k, v)
				first = false
			} else {
				fmt.Fprintf(p, "&%s=%v", k, v)
			}
		}
		dsn = p.String()
	}

	return dsn
}

type mysqldParams struct {
	protocol string
	host string
	user string
	passwd string
	db string
	attrs map[string]interface{}
}

type mysqlOption func(*mysqldParams)

func DomainSocket(sock string) mysqlOption {
	return func(params *mysqldParams) {
		params.protocol = "unix"
		params.host = sock
	}
}
func Host(host string, port ...int) mysqlOption {
	return func(params *mysqldParams) {
		params.protocol = "tcp"

		h := "localhost"
		if len(host) > 0 {
			h = host
		}

		p := 3306
		if len(port) > 0 && port[0] > 0 {
			p = port[0]
		}

		params.host = fmt.Sprintf("%s:%d", h, p)
	}
}
func User(user string, passwd ...string) mysqlOption {
	return func(params *mysqldParams) {
		params.user = user
		if len(passwd) > 0 {
			params.passwd = passwd[0]
		}
	}
}
func DBName(db string) mysqlOption {
	return func(params *mysqldParams) {
		params.db = db
	}
}
func Attr(attr string, val interface{}) mysqlOption {
	return func(params *mysqldParams) {
		params.attrs[attr] = val
	}
}
func Attrs(attrs map[string]interface{}) mysqlOption {
	return func(params *mysqldParams) {
		for attr, val := range attrs {
			params.attrs[attr] = val
		}
	}
}
