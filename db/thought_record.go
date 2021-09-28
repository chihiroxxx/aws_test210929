package db

import (
	"database/sql"
	"fmt"
	"sort"
)

func Thought_record_create(t *Thoughts) error {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't create... %v", err)
	}
	defer db.Close()

	cmd := `INSERT INTO thoughts
					(idea, page, reading_time, date, book_id)
					VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(cmd, t.Idea, t.Page, t.Reading_time, t.Date, t.Book_id)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't thought create... %v", err)
	}
	fmt.Println("thought record作りました！！！")
	return nil
}

// func Thought_record_index(jt *JoinedThoughtRecord) ([]JoinedThoughtRecord, error) {
func Thought_record_index(userid string) ([]JoinedThoughtRecord, error) {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("can't mysql open... %v", err)
	}
	defer db.Close()

	cmd := `SELECT thoughts.id, thoughts.idea, thoughts.page, thoughts.reading_time, thoughts.date, books.booktitle ,
	 					books.author, books.bookimage, books.created_at,books.updated_at, books.id ,books.user_id
	 					FROM thoughts
						JOIN books
						ON thoughts.book_id = books.id
						WHERE books.user_id = ?` // ただこれだと全くもって全件取り出しちゃうんだよなあ
	// つまり他のユーザーのものも見れてしまう...
	// WHERE books.user_id = 1 JOINでもWHEREできるっぽい！！！

	rows, err := db.Query(cmd, userid)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("can't get index... %v", err)
	}
	defer rows.Close()
	var index []JoinedThoughtRecord
	for rows.Next() {
		var jt JoinedThoughtRecord
		err := rows.Scan(&jt.Id, &jt.Idea, &jt.Page,
			&jt.Reading_time, &jt.Date, &jt.Booktitle,
			&jt.Author, &jt.Bookimage, &jt.Created_at, &jt.Updated_at,
			&jt.Book_id, &jt.User_id,
		)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("can't get row next on mysql... %v", err)
		}
		index = append(index, jt)
	}
	// ここで、並び替えしよう！！（Thought ID順でいいかな？？）
	// 他の順番のがいいかな？？？？？
	sort.Slice(index, func(i, j int) bool { return index[i].Id < index[j].Id }) // これでThought ID が小さい順になったけど...

	// fmt.Println(index)
	return index, nil
}

type JoinedThoughtRecord struct {
	// thoughtのcolumnたち
	Id   int    `json:"id"`
	Idea string `json:"thoughts"` //reactがエラーになるから...dateにする...
	// Idea         string `json:"idea"` //reactがエラーになるから...dateにする...
	Page         int    `json:"page"`
	Reading_time int    `json:"readingtime"` //ストップウォッチ的な差分もOKなように、純粋なint
	Date         string `json:"date"`
	// Created_at   string `json:"date"`        //reactがエラーになるから...dateにする...
	Created_at string `json:"createdat"` //修正中！！！(date追加中...)
	Updated_at string `json:"updatedat"`
	Book_id    int    `json:"bookid"`
	// bookのcolumnたち
	Booktitle string `json:"booktitle"`
	Author    string `json:"author"`
	Bookimage string `json:"bookimage"`
	User_id   int    `json:"userid"`
}

func Thought_record_delete(id string) error {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't delete... %v", err)
	}
	defer db.Close()
	/*
		i, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
		}
	*/
	fmt.Println("idは", id)
	// うーん id がなくて更新できなかった場合でも、反応が変わらないから...
	// 分岐させないとなああ
	cmd := `DELETE FROM thoughts WHERE id = ?`
	_, err = db.Exec(cmd, id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't delete... %v", err)
	}
	return nil
}

func Thought_record_update(t *Thoughts) error {
	db, err := sql.Open(SQL_DRIVER, SQL_CONFIG+dbname)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't mysql open... %v", err)
	}
	defer db.Close()
	// うーん id がなくて更新できなかった場合でも、反応が変わらないから...
	// 分岐させないとなああ
	cmd := `UPDATE thoughts SET idea = ? WHERE id = ?`
	// cmd := `UPDATE books SET thoughts = ? WHERE id = ?` // 注意！！column名！！！
	// _, err = db.Exec(cmd, "updated!!! form golang", 3)
	_, err = db.Exec(cmd, t.Idea, t.Id)
	// _, err = db.Exec(cmd, b.Thoughts, b.Id) // 注意！field名！！！
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("can't update... %v", err)
	}
	return nil
}
