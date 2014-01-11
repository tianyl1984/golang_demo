package dbDemo

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//访问数据库
/*
CREATE TABLE `userinfo` (
	`id` INT(10) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(50) NULL DEFAULT NULL,
	`age` INT NULL DEFAULT NULL,
	`email` VARCHAR(50) NULL DEFAULT NULL,
	PRIMARY KEY (`id`)
)
*/

func StartDemo() {
	db, err := sql.Open("mysql", "root:tyl6632460@tcp(localhost:3306)/mytest1?charset=utf8")
	checkError(err)
	//insert(db)
	//update(db)
	//remove(db)
	find(db)
}

func remove(db *sql.DB) {
	fmt.Println("<---------start delete--------->")
	stmt, err := db.Prepare("delete from userinfo where id = ?")
	checkError(err)
	res, err := stmt.Exec(4)
	rows, err := res.RowsAffected()
	fmt.Println("delete:", rows)
	fmt.Println("<---------start delete--------->")
}

func insert(db *sql.DB) {
	fmt.Println("<---------start insert--------->")
	stmt, err := db.Prepare("insert into userinfo(name,age,email) value(?,?,?)")
	checkError(err)
	res, err := stmt.Exec("zhangsan", 128, "zhangsan@126.com")
	checkError(err)
	id, err := res.LastInsertId()
	checkError(err)
	fmt.Println("insert:", id)
	fmt.Println("<---------end insert----------->")
}

func update(db *sql.DB) {
	fmt.Println("<---------start update--------->")
	stmt, err := db.Prepare("update userinfo set email = ? where id = ?")
	checkError(err)
	rs, err := stmt.Exec("lisi@163.com", 4)
	checkError(err)
	fmt.Println(rs.RowsAffected())
	fmt.Println("<---------end update----------->")
}

func find(db *sql.DB) {
	fmt.Println("<---------start select--------->")
	rows, err := db.Query("select * from userinfo")
	for rows.Next() {
		var id int
		var name string
		var age int
		var email string
		err = rows.Scan(&id, &name, &age, &email) //有null值报错
		checkError(err)
		fmt.Printf("ID:%d,name:%v,age:%d,email:%v\r\n", id, name, age, email)
	}
	fmt.Println("<---------end select----------->")
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
