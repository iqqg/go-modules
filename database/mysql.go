package database

import (
	"database/sql"
	"fmt"
	"log"

	// just for init
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db1, err := sql.Open("mysql", "pi:shine@tcp(192.168.1.4)/test")
	if err != nil {
		panic(err.Error())
	}
	db = db1
}

func query(str string) *sql.Rows {
	stmt, err := db.Prepare(str)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return rows
}

func mysqlExample() {
	defer db.Close()

	stmtIns, err := db.Prepare("insert into person (name, age, pdesc) values (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	for i := 0; i < 30; i++ {
		_, err = stmtIns.Exec(fmt.Sprintf("Albert%d", i), i, "一个普通人")
		if err != nil {
			panic(err.Error())
		}
	}

	rows := query("select name, age, pdesc from person")
	if rows == nil {
		panic(false)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name string
			age  int
			desc string
		)
		if err = rows.Scan(&name, &age, &desc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("name: %v age: %d desc:%v\n", name, age, desc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
