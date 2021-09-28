package web

import (
	"backend/db"
	"fmt"
	// "encoding/csv"
	// "io/ioutil"
	// "os"

	// "log"
	"net/http"

	// "net/url"
	"strconv"

	"github.com/labstack/echo"
)

// csv はのち！！ JOIN INDEXから持ってくる！！！
/*
func createCsv(userid string) error {
	file, err := os.Create("./csv/index.csv")
	if err != nil {
		fmt.Println(err)
		// return c.JSON(http.StatusBadRequest, "no record...")
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	w.Write([]string{"No", "BookTitle", "Author", "BookImage", "Thoughts", "Date", "CreatedAt", "UpdatedAt"})

	index, err := db.Book_record_index(userid)
	if err != nil {
		fmt.Println(err)
		// return c.JSON(http.StatusBadRequest, "no record...")
		return err
	}
	// index[0].Id こんな風にアクセスするのか...
	for i := 0; i < len(index); i++ {
		var arr []string
		arr = append(arr, strconv.Itoa(i+1), index[i].Booktitle, index[i].Author,
			index[i].Bookimage, index[i].Thoughts, index[i].Date,
			index[i].Created_at, index[i].Updated_at,
			// , strconv.Itoa(index[i].User_id)
		)
		w.Write(arr)
	}

	// return c.Inline("./csv/index.csv", "index.csv")
	return nil
}

func BookCsv(c echo.Context) error {
	userid := c.Param("userid")
	err := createCsv(userid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, "no record...")
	}
	byteFile, err := ioutil.ReadFile("./csv/index.csv")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, "no record...")
	}
	// test := []byte(`aaa,test,aaa,ts`)
	return c.Blob(http.StatusOK, "text/csv", byteFile)
}
*/

//現在はこのURLにアクセスはなし！！ コレクションページの時に使う予定！！
/*
func BookIndex(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(c)
	index, err := db.Book_record_index(id)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "no record...")
	}
	// thoughtとJOINのテスト！！！（下記web_thoughtに移行済み！）
	// index, err := db.Thought_record_index(id) //結局渡すのはuseridで
	// if err != nil {
	// 	fmt.Println(err)
	// 	return echo.NewHTTPError(http.StatusBadRequest, "no record...")
	// }

	return c.JSON(http.StatusOK, index) //{books: index}で返却したいんだけども...
}
*/
func BookFinish(c echo.Context) error {
	id := c.Param("id") //book_idだよ！
	err := db.Book_record_finish_count_up(id)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't count up...")
	}
	return c.JSON(http.StatusOK, "finish count up OK!")
}

func BookIndex(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(c)
	index, err := db.Book_record_index(id)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "no record...")
	}
	// thoughtとJOINのテスト！！！（下記web_thoughtに移行済み！）
	// index, err := db.Thought_record_index(id) //結局渡すのはuseridで
	// if err != nil {
	// 	fmt.Println(err)
	// 	return echo.NewHTTPError(http.StatusBadRequest, "no record...")
	// }

	return c.JSON(http.StatusOK, index) //{books: index}で返却したいんだけども...
}

/*
func BookUpdate(c echo.Context) error {
	b := updateRecord(c)
	fmt.Println(b)
	db.Book_record_update(b)
	return c.JSON(http.StatusOK, "book updated!!!")
}
*/
/*
func BookDelete(c echo.Context) error {
	fmt.Println(c.Request().URL.String())
	// b := deleteRecord(c)
	// fmt.Println(b)
	// id := c.QueryParam("id")
	id := c.Param("id")
	fmt.Println(id)
	err := db.Book_record_delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "book deleted!!!")
}
*/
func BookCreate(c echo.Context) error {
	/*
		// paramsをqueryでとってくる
		// q := c.QueryParam("q") //これは一つのquery
		q := c.QueryParams()
		fmt.Println(q)
		fmt.Println("qは" + q.Get("q"))
	*/

	// p := c.FormValue("booktitle")
	// // p, _ := c.FormParams()
	// fmt.Println(p)

	// バリデーション検証（ストロングパラメーター）

	// ここで、Bookあるなし判定する？？？
	b := newRecord(c)
	bookexistid, err := bookExistChecker(b) // bookidとerrを返すから...
	// でなければ、一応-1にしてる つまり-1ならbookが今まで存在していないといえそう。
	fmt.Println(bookexistid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	// fmt.Println(b)
	// if のバリデーション通してOKだったら下の処理
	// で、DBに登録
	if bookexistid == -1 {
		createdid, err := db.Book_record_create(b)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		fmt.Println(createdid)
		t := newThoughtRecord(c, createdid) //ちがう、大本を呼び出して作るだけにしよう。
		db.Thought_record_create(t)

	} else {
		// で、ここまでで、今までbook登録しているのか、判定できたから、
		// thoughtを作成する関数にバトンタッチしよう！！
		// dbでいいのか、

		t := newThoughtRecord(c, bookexistid) //ちがう、大本を呼び出して作るだけにしよう。
		db.Thought_record_create(t)
		//うーーーん...ここまでthoughtだけど...

	}

	return c.JSON(http.StatusOK, "book created!!!")
}

func bookExistChecker(b *db.Book) (int, error) {
	// まっっっっっって！！！これじゃ他のユーザーがいたら...うんんうんかんぬmん...
	bookexistid, err := db.Book_record_exist_checker(b)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return bookexistid, nil
}

// func newRecord(q url.Values) *db.Book {
func newRecord(c echo.Context) *db.Book { // 旧Bookテーブルだから！！ね！
	// swiftからのformvalue test!!!
	fmt.Println(c.FormValue("booktitle"))

	// 注意すること！！
	b := new(db.Book)
	// 色々処理------------------
	b.Booktitle = c.FormValue("booktitle")
	b.Author = c.FormValue("author")
	b.Bookimage = c.FormValue("bookimage")
	/*
		b.Thoughts = c.FormValue("thoughts")
		if c.FormValue("page") != "" {
			page, _ := strconv.Atoi(c.FormValue("page"))
			b.Page = page
		}
		if c.FormValue("readingtime") != "" {
			readingtime, _ := strconv.Atoi(c.FormValue("readingtime"))
			b.Reading_time = readingtime
		}
	*/
	// b.Date = c.FormValue("date")
	user_id, err := strconv.Atoi(c.FormValue("userid"))
	// err ハンドリングするために、strconvを他に代入してからにするか
	if err != nil {
		fmt.Println(err)
	}
	b.User_id = user_id
	return b
}

// 今のこのメソッド使ってはないよ...(そうか、現状bookの情報をUPDATEすることはないからか、 count upはするけども...)
func updateRecord(c echo.Context) *db.Book {
	b := new(db.Book)
	// 色々処理------------------
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
	}
	b.Id = id

	// b.Thoughts = c.FormValue("thoughts")
	if c.FormValue("date") != "" {
		// b.Date = c.FormValue("date")
	}
	/*
		user_id, err := strconv.Atoi(c.FormValue("userid"))
		// err ハンドリングするために、strconvを他に代入してからにするか
		if err != nil {
			fmt.Println(err)
		}
		b.User_id = user_id
	*/
	return b
}

/*
func deleteRecord(c echo.Context) *db.Book {
	b := new(db.Book)
	// 色々処理------------------
	id, err := strconv.Atoi(c.FormValue("id")) //deleteはformdataじゃないのか？？？

	if err != nil {
		fmt.Println(err)
	}
	b.Id = id
	return b
}
*/
