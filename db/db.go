package db

import (
	"database/sql"
	"fmt"
	"os"

	// "os"

	_ "github.com/go-sql-driver/mysql"
)

// var DbConnection *sql.DB
// var name string

/*
var err error
func init() {
	db_create()
	db, err = sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

}
*/
func Db_main() {
	// db, err := sql.Open("mysql", "root:(ここはパスワードみたい！！！)@/test_database（使うデータベース名！）")
	// db, err := sql.Open("mysql", "root:@/backRakuten_development")
	// db, err := sql.Open("mysql", "root:@/test_go_test")

	db_create()

	// user_record_create() //user作る

	// book_record_create() //book登録 <---echoで呼び出す！

	/*

		// err = db.QueryRow("SELECT booktitle FROM books WHERE id = ?", 3).Scan(&name)
		rows, err := db.Query("SELECT * FROM books")
		if err != nil {
			fmt.Println(err)
		}

		defer rows.Close()

		for rows.Next() {
			var book Book
			err := rows.Scan(&book.Id, &book.Booktitle, &book.Author, &book.Bookimage, &book.Thoughts, &book.Date, &book.Created_at,
				&book.Updated_at, &book.User_id)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(book)
		}


	*/
	fmt.Println(sql.Drivers()) // mysql と表示される
}

type Book struct {
	Id        int    `json:"id"`
	Booktitle string `json:"booktitle"`
	Author    string `json:"author"`
	Bookimage string `json:"bookimage"`
	Bookurl   string `json:"bookurl"` //追加！！！
	// Thoughts     string `json:"thoughts"`
	// Page         int    `json:"page"`
	// Reading_time int    `json:"readingtime"` //ストップウォッチ的な差分もOKなように、純粋なint
	Finish_count int `json:"finishcount"` //ストップウォッチ的な差分もOKなように、純粋なint
	// Date         string `json:"date"`
	Created_at string `json:"createdat"`
	Updated_at string `json:"updatedat"`
	User_id    int    `json:"userid"`
}

type User struct { //email追加したいなあ...
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Password_digest string `json:"passworddigest"`
	Created_at      string `json:"createdat"`
	Updated_at      string `json:"updatedat"`
}

type Thoughts struct {
	Id           int    `json:"id"`
	Idea         string `json:"idea"`
	Page         int    `json:"page"`
	Reading_time int    `json:"readingtime"` //ストップウォッチ的な差分もOKなように、純粋なint
	Date         string `json:"date"`
	Created_at   string `json:"createdat"`
	Updated_at   string `json:"updatedat"`
	Book_id      int    `json:"bookid"`
}

type Finish struct {
	Id          int    `json:"id"`
	Count       int    `json:"count"`
	Finish_date string `json:"finishdate"`
	Created_at  string `json:"createdat"`
	Updated_at  string `json:"updatedat"`
	Book_id     int    `json:"bookid"`
}

var db *sql.DB                //グローバルアクセスのため
var dbname = "test_go_test10" //これもレコード登録時に使う

const SQL_DRIVER string = "mysql" // ん？mysqlはそのままでいいか。 ENVにする！！
// const SQL_CONFIG string = "root:@/" // ENVにする！！
// const SQL_CONFIG string = "root:password@tcp(memento_mysql:3306)/" // Docker用！！
var DATABASE_USERNAME string = os.Getenv("DATABASE_USERNAME")
var DATABASE_PASSWORD string = os.Getenv("DATABASE_PASSWORD")
var DATABASE_HOSTNAME string = os.Getenv("DATABASE_HOSTNAME") //rdsのエンドポイントでいいとのこと！！！

// const SQL_CONFIG string = "root:password@tcp(memento_mysql:3306)/" // Docker用！！
var SQL_CONFIG string = DATABASE_USERNAME + ":" + DATABASE_PASSWORD + "@tcp(" + DATABASE_HOSTNAME + ":3306)/" // Docker用！！

func db_create() {
	// DataSourceName user:password@tcp(container-name:port)/dbname
	// DataSourceName = "root:golang@tcp(mysql-container:3306)/golang_db"
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG) // これでDB未選択で接続できてるはず！！！
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	//カラム検証テストのためのDB削除文！！！-------------------
	/*
		cmd := `DROP DATABASE ` + dbname
		_, err = db.Exec(cmd)
		if err != nil {
			fmt.Println(err)
		}
	*/
	// ここまで--------------------------------------------

	cmd := `CREATE DATABASE IF NOT EXISTS ` + dbname // これで test_go_test2 っていうDBはできた！！ 変数化もできた！
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	}

	cmd = `USE ` + dbname
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	}

	// column にpageも持たせたいなあ...
	// emailも追加したい！！
	// あとなんだろう？？
	// user_id BIGINT のreference は子テーブルで作るんだって！

	// users テーブル！！
	cmd = `CREATE TABLE IF NOT EXISTS users (id BIGINT auto_increment primary key,
		name VARCHAR(255) unique NOT NULL,
		password_digest VARCHAR(255) NOT NULL,
		created_at DATETIME default current_timestamp NOT NULL,
		updated_at DATETIME default current_timestamp on update current_timestamp NOT NULL
		)`
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	}

	// books テーブル！ usersテーブルが親
	cmd = `CREATE TABLE IF NOT EXISTS books
	(id BIGINT auto_increment primary key,
		booktitle VARCHAR(255),
		author VARCHAR(255),
		bookimage VARCHAR(255),
		bookurl VARCHAR(255),
		finish_count BIGINT default 0,
		created_at DATETIME default current_timestamp NOT NULL,
		updated_at DATETIME default current_timestamp on update current_timestamp NOT NULL,
		user_id BIGINT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
		)`
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	}
	// thoughts テーブル！ booksテーブルが親
	cmd = `CREATE TABLE IF NOT EXISTS thoughts (id BIGINT auto_increment primary key,
		idea VARCHAR(255),
		page BIGINT default 0,
		reading_time BIGINT default 0,
		date DATETIME default current_timestamp,
		created_at DATETIME default current_timestamp NOT NULL,
		updated_at DATETIME default current_timestamp on update current_timestamp NOT NULL,
		book_id BIGINT NOT NULL,
		FOREIGN KEY(book_id) REFERENCES books(id)
		)`
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	}

	// finishes テーブル！ booksテーブルが親
	cmd = `CREATE TABLE IF NOT EXISTS finishes (id BIGINT auto_increment primary key,
		count BIGINT default 0,
		finish_date DATETIME default current_timestamp,
		created_at DATETIME default current_timestamp NOT NULL,
		updated_at DATETIME default current_timestamp on update current_timestamp NOT NULL,
		book_id BIGINT NOT NULL,
		FOREIGN KEY(book_id) REFERENCES books(id)
		)`
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	}
}
