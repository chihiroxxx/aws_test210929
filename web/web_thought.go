package web

import (
	"backend/db"
	"encoding/csv"
	"io/ioutil"
	"os"

	// "encoding/json"
	"fmt"
	"net/http"
	"strconv"

	// set "github.com/deckarep/golang-set"
	"github.com/labstack/echo"
)

// func ThoughtCreate(c echo.Context) error { //ん？ここでintも引数にしちゃうと
// 	// 後々にhttpから呼べなくなりそう...

// 	newThoughtRecord(c)
// }
func newThoughtRecord(c echo.Context, bookexistid int) *db.Thoughts { // 旧Bookテーブルだから！！ね！
	// 注意すること！！
	t := new(db.Thoughts)
	// 色々処理------------------
	t.Idea = c.FormValue("thoughts")
	if c.FormValue("page") != "" {
		page, _ := strconv.Atoi(c.FormValue("page"))
		t.Page = page
	}
	if c.FormValue("readingtime") != "" {
		readingtime, _ := strconv.Atoi(c.FormValue("readingtime"))
		t.Reading_time = readingtime
	}
	t.Date = c.FormValue("date")
	// book_id, err := strconv.Atoi(c.FormValue("userid"))
	// // err ハンドリングするために、strconvを他に代入してからにするか
	// if err != nil {
	// 	fmt.Println(err)
	// }
	t.Book_id = bookexistid
	return t
}

func ThoughtIndex(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(c)
	/* thoughtとJOINのテスト！！！
	index, err := db.Book_record_index(id)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "no record...")
	}
	*/
	index, err := db.Thought_record_index(id) //結局渡すのはuseridで
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "no record...")
	}
	return c.JSON(http.StatusOK, index) //{books: index}で返却したいんだけども...
}

func ThoughtDelete(c echo.Context) error {
	// fmt.Println(c.Request().URL.String())
	// b := deleteRecord(c)
	// fmt.Println(b)
	// id := c.QueryParam("id")
	id := c.Param("id") //ここは今までbooks.idだったのか
	fmt.Println(id)
	err := db.Thought_record_delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "thought deleted!!!")
}

func ThoughtUpdate(c echo.Context) error {
	t := updateThoughtRecord(c)
	err := db.Thought_record_update(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "thought updated!!!")
}

func updateThoughtRecord(c echo.Context) *db.Thoughts {
	t := new(db.Thoughts)
	// 色々処理------------------
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
	}
	t.Id = id

	// b.Thoughts = c.FormValue("thoughts")
	t.Idea = c.FormValue("thoughts") // b.Thoughts -> t.Ideaなの注意...
	//つまり現状へんこうできるのはthoughtsだけ？
	// page,readingtimeもへんこうさせる？？？
	if c.FormValue("page") != "" { //まだreactでは実装させてないから ずっと ==""
		t.Page, err = strconv.Atoi(c.FormValue("page"))
		if err != nil {
			fmt.Println(err)
		}
	}

	if c.FormValue("readingtime") != "" { //まだreactでは実装させてないから ずっと ==""
		t.Reading_time, err = strconv.Atoi(c.FormValue("readingtime"))
		if err != nil {
			fmt.Println(err)
		}
	}

	/*
		if c.FormValue("date") != "" {
			b.Date = c.FormValue("date") // column作らないと処理できなし
		}
	*/
	/*
		user_id, err := strconv.Atoi(c.FormValue("userid"))
		// err ハンドリングするために、strconvを他に代入してからにするか
		if err != nil {
			fmt.Println(err)
		}
		b.User_id = user_id
	*/
	return t
}

// csv はのち！！ JOIN INDEXから持ってくる！！！

func createThoughtCsv(userid string) error {
	file, err := os.Create("./csv/index.csv")
	if err != nil {
		fmt.Println(err)
		// return c.JSON(http.StatusBadRequest, "no record...")
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	w.Write([]string{"No", "BookTitle", "Author", "BookImage", "Idea", "Date", "CreatedAt", "UpdatedAt"})

	index, err := db.Thought_record_index(userid)
	if err != nil {
		fmt.Println(err)
		// return c.JSON(http.StatusBadRequest, "no record...")
		return err
	}
	// index[0].Id こんな風にアクセスするのか...
	for i := 0; i < len(index); i++ {
		var arr []string
		arr = append(arr, strconv.Itoa(i+1), index[i].Booktitle, index[i].Author,
			index[i].Bookimage, index[i].Idea, index[i].Date,
			index[i].Created_at, index[i].Updated_at,
			// , strconv.Itoa(index[i].User_id)
		)
		w.Write(arr)
	}

	// return c.Inline("./csv/index.csv", "index.csv")
	return nil
}

func ThoughtCsv(c echo.Context) error {
	userid := c.Param("userid")
	err := createThoughtCsv(userid)
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
