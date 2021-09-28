package db

import (
	"database/sql"
	"fmt"
	// "log"
	// "strconv"
)

/*
func book_record_init() {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

}
*/

func Book_record_exist_checker(b *Book) (int, error) { // intじゃなくてstringで..のがいいのか？？
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("can't mysql open ... %v", err)
	}
	defer db.Close()

	if b.Booktitle == "" {
		fmt.Println(err)
		return -1, fmt.Errorf("can't search blanc title... %v", err)
	}

	//ようはここで、自分のuserid & booktitleにマッチするか調べてるってこと
	cmd := `SELECT * FROM books WHERE booktitle = ? AND user_id = ?`
	rows, err := db.Query(cmd, b.Booktitle, b.User_id)
	if err != nil {
		fmt.Println(err)
		// return -1, fmt.Errorf("can't get book , it don't exist %v", err)
		return -1, nil //これはマッチするbooktitleがないですよーってことでしょ？
		// あ、ここでエラーにしたらダメじゃないの...
	}
	defer rows.Close()
	// --------------------------------------------------------------

	var index []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.Id, &b.Booktitle, &b.Author, &b.Bookimage,
			&b.Bookurl, &b.Finish_count,
			// &b.Thoughts, &b.Page, &b.Reading_time,
			&b.Created_at,
			&b.Updated_at, &b.User_id)
		if err != nil {
			fmt.Println(err)
			return -1, fmt.Errorf("can't get row next... %v", err)
		}
		index = append(index, b)
	}
	for _, p := range index {
		fmt.Println(p.Booktitle)
	}

	// return index[0].Id, nil //ここでエラーっぽい。
	// fmt.Println(index[0].Id) lengthだ！！！！！
	fmt.Println("indexのlengthは...？->", len(index))
	if len(index) == 0 {
		return -1, nil //これはマッチするbooktitleがないですよーってことでしょ？
	}
	return index[0].Id, nil //ここでエラーっぽい。テストで
	// テストで   ベタ書き！！！！！！！！！！！！！

	// そうか、もし無かった何が返ってくるのかなって思ったんだけど...
	// そもそも、その前に-1 と エラーで早期リターんしてくれるのか
}

func Book_record_create(b *Book) (int, error) {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("can't mysql open... %v", err)
	}
	defer db.Close()

	// cmd := `INSERT INTO books // テーブル分割につき...
	// 				(booktitle, author, bookimage, thoughts, page, reading_time, date, user_id)
	// 				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	// result, err := db.Exec(cmd, b.Booktitle, b.Author, b.Bookimage, b.Thoughts, b.Page, b.Reading_time, b.Date, b.User_id)
	cmd := `INSERT INTO books
					(booktitle, author, bookimage,
						bookurl, finish_count, user_id)
					VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(cmd, b.Booktitle, b.Author, b.Bookimage, b.Bookurl, b.Finish_count, b.User_id)

	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("can't create... %v", err)
	}
	fmt.Println("book record作りました！！！")
	fmt.Println(result.LastInsertId())
	createdid, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)

	}
	// ここでfinishレコードを作るってことか、
	err = Finish_record_first_create(int(createdid))
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("can't create finishrecord... %v", err)
	}

	return int(createdid), nil ///そりゃ...そうだ...このbは送られてきたb.Idじゃんよ
	// つまり作ったレコードIdを入れてるわけじゃないぜえええええ

}

/*
func Book_record_update(b *Book) {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	// うーん id がなくて更新できなかった場合でも、反応が変わらないから...
	// 分岐させないとなああ
	cmd := `UPDATE books SET thoughts = ? WHERE id = ?`
	// _, err = db.Exec(cmd, "updated!!! form golang", 3)
	_, err = db.Exec(cmd, b.Thoughts, b.Id)
	if err != nil {
		fmt.Println(err)

	}
}
*/
/*
func Book_record_delete(id string) error {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't delete... %v", err)
	}
	defer db.Close()

		// i, err := strconv.Atoi(id)
		// if err != nil {
		// 	fmt.Println(err)
		// }

	fmt.Println("idは", id)
	// うーん id がなくて更新できなかった場合でも、反応が変わらないから...
	// 分岐させないとなああ
	cmd := `DELETE FROM books WHERE id = ?`
	_, err = db.Exec(cmd, id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't delete... %v", err)
	}
	return nil
}
*/
/*
func Book_record_index(id string) ([]Book, error) {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("can't open mysql %v", err)
	}
	defer db.Close()

	fmt.Println("indexのidは" + id)
	if id == "" {
		fmt.Println(err)
		return nil, fmt.Errorf("can't get any id... %v", err)
	}
	cmd := `SELECT * FROM books WHERE user_id = ?`
	rows, err := db.Query(cmd, id)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("can't get index... %v", err)
	}
	defer rows.Close()

	var index []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.Id, &b.Booktitle, &b.Author, &b.Bookimage,
			&b.Thoughts, &b.Page, &b.Reading_time,
			&b.Date, &b.Created_at,
			&b.Updated_at, &b.User_id)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("can't get row next... %v", err)
		}
		index = append(index, b)
	}
	for _, p := range index {
		fmt.Println(p.Booktitle)
	}
	return index, nil

}
*/

// Bookのcollectionを返却する！！！
// 次！！！ここの編集！！！
// func Book_record_finish_count_up(id string) error {
// 	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
// 	if err != nil {
// 		fmt.Println(err)
// 		return fmt.Errorf("can't open mysql %v", err)
// 	}
// 	defer db.Close()
// 	cmd := `UPDATE books SET finish_count = finish_count + 1 WHERE id = ?`
// 	_, err = db.Exec(cmd, id)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return nil
// }
// 次！！！ここの編集！！！
func Book_record_finish_count_up(id string) error { //このidはbooksのレコードのidか
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't open mysql %v", err)
	}
	defer db.Close()
	// ここがこんどUPDATEじゃなくてINSERT INTO finishes うんぬんカンヌン....+1
	// 今度は+1が問題か...
	// cmd := `UPDATE books SET finish_count = finish_count + 1 WHERE id = ?`
	// そう、要はこの count のところに入れるvalueが問題なのだ...
	// これもヒント SELECT MAX(page) FROM thoughts WHERE book_id = 2;
	cmd := `INSERT INTO finishes (count, book_id) SELECT (SELECT MAX(count) FROM finishes WHERE book_id = ?) + 1, ?;`
	_, err = db.Exec(cmd, id, id)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func Book_record_index(id string) ([]Book, error) {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("can't open mysql %v", err)
	}
	defer db.Close()

	fmt.Println("indexのidは" + id)
	if id == "" {
		fmt.Println(err)
		return nil, fmt.Errorf("can't get any id... %v", err)
	} else if id == "-1" {
		fmt.Println(err)
		return nil, fmt.Errorf("can't get -1 id... %v", err)
	}
	// ここの取得 Queryが変わるのか...むむ...
	// cmd := `SELECT * FROM books WHERE user_id = ?`
	// cmd := `SELECT books.id, books.booktitle, books.author, books.bookimage,
	// 								books.bookurl, MAX(finishes.count), books.created_at,
	// 								books.updated_at, books.user_id
	// 								FROM books
	// 								JOIN finishes
	// 								ON books.id = finishes.book_id
	// 								WHERE user_id = ?` //これをどう変えればいい？
	// この MAX(finishes.count)が故に、レコード一件しか取ってこれないのか....？（そうみたい 検証済み）
	// cmd := `SELECT books.id, books.booktitle, books.author, books.bookimage,
	// 								books.bookurl, finishes.count, books.created_at,
	// 								books.updated_at, books.user_id
	// 								FROM books
	// 								JOIN finishes
	// 								ON books.id = finishes.book_id
	// 								WHERE user_id = ?` //これをどう変えればいい？
	// これか！？！？
	cmd := `SELECT
									booksA.id, booksA.booktitle,
									booksA.author, booksA.bookimage,
									booksA.bookurl, finishesA.count, booksA.created_at,
	 								booksA.updated_at, booksA.user_id
						FROM
							(books AS booksA,
							finishes AS finishesA)
							INNER JOIN (SELECT
															book_id,
															MAX(count) AS MaxCount
													FROM
															finishes
													GROUP BY
															book_id) AS finishesB
							ON booksA.id = finishesB.book_id
							AND finishesA.count = finishesB.MaxCount
							WHERE user_id = ?
							GROUP BY id`
	// WHERE finishes.count = ( SELECT MAX(count) FROM finishes)
	// AND user_id = ?` //これをどう変えればいい？
	// SELECT MAX(page) FROM thoughts WHERE book_id = 2;
	//  cmd := `SELECT thoughts.id, thoughts.idea, thoughts.page, thoughts.reading_time, thoughts.date, books.booktitle ,
	//  					books.author, books.bookimage, books.created_at, books.id ,books.user_id
	//  					FROM thoughts
	// 					JOIN books
	// 					ON thoughts.book_id = books.id
	// 					WHERE books.user_id = ?`
	rows, err := db.Query(cmd, id)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("can't get index... %v", err)
	}
	defer rows.Close()

	var index []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.Id, &b.Booktitle, &b.Author, &b.Bookimage,
			&b.Bookurl, &b.Finish_count,
			// &b.Date,
			&b.Created_at,
			&b.Updated_at, &b.User_id)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("can't get row next... %v", err)
		}
		index = append(index, b)
	}
	for _, p := range index {
		fmt.Println(p.Booktitle)
	}
	return index, nil

}
