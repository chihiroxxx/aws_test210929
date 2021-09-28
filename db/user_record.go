package db

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// このファイルは必ずユーザーレコードに関することにするから、
// initでusers tableを指定するなどをまとめる。
// 一旦init置いといてシンプルに作ってみる！
func user_record_init() {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

}

func User_record_delete(u *User) error {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("don't delete... %v", err)
	}
	defer db.Close()

	cmd := `DELETE FROM users WHERE name = ?`
	_, err = db.Exec(cmd, u.Name)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("don't delete... %v", err)
	}
	return nil
}

// func User_record_checker(u *User) error {
// 	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer db.Close()

// }

func User_record_create(u *User) (int, error) {
	// user_record_init()
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	cmd := `INSERT INTO users (name, password_digest) VALUES (?, ?)`
	result, err := db.Exec(cmd, u.Name, u.Password_digest) /* このpasswordをハッシュ化して->保存の作業の
	ハッシュ化のところはdbでやる作業じゃないよね？？？
	バリデーションかけて、OKだったら、ハッシュ化して
	DBに保存 だよね！！
	つまり、この保存の機構はこれでOK！ */
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("can't create user... %v", err)
	}
	fmt.Println("Userのrecord作りました！！！")
	resultId, _ := result.LastInsertId()
	fmt.Println(resultId)
	// int64をint型に コンバートする！！！
	// そんでreturnする！！
	return int(resultId), nil
}

func User_record_login(u *User) (int, error) {

	// user_record_init()
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	cmd := `SELECT * FROM users WHERE name = ?`

	rows, err := db.Query(cmd, u.Name)
	defer rows.Close()

	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id,
			&user.Name,
			&user.Password_digest,
			&user.Created_at,
			&user.Updated_at,
		)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
	}

	fmt.Println(user.Password_digest)
	if err != nil {
		fmt.Println(err)
	}

	// DBに保存してある、passwordと reactのログインフォームから送られてきた、passwordを比較！
	// で、一致しなければエラーを返す！！   一致していればnil つまり何も返さない！
	err = bcrypt.CompareHashAndPassword([]byte(user.Password_digest), []byte(u.Password_digest))
	if err != nil {
		// fmt.Println(err) // あああああ！！！ log.Fatalln が強制終了させるのか！？！？！？
		// Fatalln is equivalent to Println() followed by a call to os.Exit(1). だって
		// あああああああああ！！！！！
		fmt.Println(err)
		return -1, err
	}
	/*
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			fmt.Println(err)
		}

		defer rows.Close()

		for rows.Next() {
			var user User
			err := rows.Scan(&user.Id,
				&user.Name,
				&user.Password_digest,
				&user.Created_at,
				&user.Updated_at,
			)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(user)
		}
	*/
	/* このpasswordをハッシュ化して->保存の作業の
	_, err = db.Exec(cmd, u.Name, u.Password_digest) このpasswordをハッシュ化して->保存の作業の
	ハッシュ化のところはdbでやる作業じゃないよね？？？
	バリデーションかけて、OKだったら、ハッシュ化して
	DBに保存 だよね！！
	つまり、この保存の機構はこれでOK！ */
	// fmt.Println(pass)

	fmt.Println("ログインできました！！！")
	return user.Id, nil
}
