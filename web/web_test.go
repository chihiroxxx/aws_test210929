package web

import (
	// "bytes"
	"backend/db"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"strings"
	"testing"

	// "io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/labstack/echo"
)

var ( // この辺りを動的するのがいいのかな？？ うーーん、テストコードの流儀？？
	testname     string = "tesMan2dsd"
	testpassword string = "12345678"
	token        string
	testuserid   string
	testbook     db.Book = db.Book{Booktitle: "test本", Author: "me", Bookimage: "httpfdsajkf.com"} // Thoughts: "ワアアアアい", //テーブル分割につき...
	// Date: "20210809123211" //テーブル分割につき...

	createdbookid string
)

var (
	base *url.URL
	err  error
	e    *echo.Echo
)

// t.RUN でサブテストなのか！！
func TestMain(m *testing.M) {
	e = echo.New()
	base, err = url.Parse("http://localhost:9090")
	if err != nil {
		fmt.Println(err)
	}
	status := m.Run()
	defer os.Exit(status)
}

func TestUserCreate(t *testing.T) {
	// e := echo.New()
	// req := new(http.Request) // httpリクエスト作成

	// r *http.Request, w http.ResponseWriter を入れて
	// echoのコンテキストを作るわけだ。
	// でそのコンテキストを testしたい メソッドに渡す！！！

	// であれば ここ共通化できそう！
	// base, err := url.Parse("http://localhost:9090")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	endUrl, err := url.Parse("/api/v1/users")
	if err != nil {
		fmt.Println(err)
	}
	endpoint := base.ResolveReference(endUrl).String()
	// fmt.Println(endpoint)

	// まずリクエストnewする
	// うーん... body, err := io.Reader.Read([]byte("testtest"),error)

	// values := url.Values{}
	values := make(url.Values)
	values.Set("name", testname)
	values.Set("password_digest", testpassword)
	//そうか...DB保存の前にバリデーションかけないとな。

	// http.MethodGet でhttpメソッド選ぶんだ！？
	// req, err := http.NewRequest(http.MethodPost,
	// 	endpoint, nil) //この...第3引数さ...やっぱりPOSTだといるのか...
	req, err := http.NewRequest(http.MethodPost,
		endpoint,
		// bytes.NewBuffer([]byte("name=tesMan")))
		strings.NewReader(values.Encode()))
	//この...第3引数さ...やっぱりPOSTだといるのか...
	//bytes.NewBuffer([]byte("name"))で第3引数に中身入れるみたい！！！
	// POSTはみられたくない情報だからbodyに入れて渡す！！
	// GETはURLのクエリパラメーターとして渡すのか！！！q=1 のやつ！
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(strings.NewReader(values.Encode()))

	// でHeaderつけるときは
	// req.Header.Set(色々Headerを書く)
	// req.Header.Addでもいい？？(色々Headerを書く)

	//このheaderが足りなかったのか！？！？！
	// echoにcontentTypeをsetしてないとなのか！？！？！
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res) // つまりechoコンテキストってさ リクエストとレスポンスライターで、できてたんだ！？！？（テストコードも面白いかも！！）

	// ここまで準備...
	// ここからいよいよテストしたいメソッドにアクセス！検証！

	err = UserCreate(c)
	// nameはuniqueなものにしているから、重複エラーでオッケ！！ でもなくない？他のerrorに気づけなくない？
	if err != nil {
		t.Error(err)
	}

	// 後処理として、createしたrecodeをDBから削除させよう...
	// むしろDBのテストも一緒に実行にすると、削除も走らせられそう。
	// // うーーーーん...
	// u := new(db.User)
	// u.Name = values.Get("name")
	// db.User_record_delete(u)
}

func TestUserLogin(t *testing.T) {
	// e := echo.New()
	// base, err := url.Parse("http://localhost:9090")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	endUrl, err := url.Parse("/api/v1/login")
	if err != nil {
		fmt.Println(err)
	}
	endpoint := base.ResolveReference(endUrl).String()

	values := make(url.Values)
	values.Set("name", testname)
	values.Set("password_digest", testpassword)

	fmt.Println(values.Encode())
	req, err := http.NewRequest(http.MethodPost,
		endpoint,
		strings.NewReader(values.Encode()))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	res := httptest.NewRecorder()
	// res := new(http.ResponseWriter)
	c := e.NewContext(req, res)
	err = UserLogin(c)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println("returnは...?", res) //あああああ！！！resに入って帰ってくるのかああああ
	// だから必要だったのかああああああああああああ！！！
	if res.Result().StatusCode != 200 {
		t.Error(err)
	}
	fmt.Println("とって来れてる？？？")
	// fmt.Println(res.Body)
	a := res.Body.Bytes()
	var result Result
	json.Unmarshal(a, &result)
	// fmt.Println(res.Result().Body.Read([]byte("token")))
	token = result.Token
	testuserid = result.Userid
}

type Result struct {
	Token  string `json:"token"`
	Userid string `json:"userid"`
}

func TestBookCreate(t *testing.T) {
	// e := echo.New()
	// base, err := url.Parse("http://localhost:9090")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	endUrl, err := url.Parse("/api/v1/restricted/books")
	if err != nil {
		fmt.Println(err)
	}

	values := make(url.Values)
	values.Set("booktitle", testbook.Booktitle)
	values.Set("author", testbook.Author)
	values.Set("bookimage", testbook.Bookimage)
	// values.Set("thoughts", testbook.Thoughts) テーブル分割につき...
	// values.Set("date", testbook.Date) //テーブル分割につき...
	values.Set("userid", testuserid)

	endpoint := base.ResolveReference(endUrl).String()
	fmt.Println(endpoint)
	req, err := http.NewRequest(http.MethodPost,
		endpoint,
		strings.NewReader(values.Encode()))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	err = BookCreate(c)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res.Body)
}

func TestBookIndex(t *testing.T) {
	// e := echo.New()

	// base, err := url.Parse("http://localhost:9090")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	endUrl, err := url.Parse("/api/v1/restricted/books/")
	if err != nil {
		fmt.Println(err)
	}
	/*
		values := make(url.Values)
		values.Add("id", "20")
	*/

	endpoint := base.ResolveReference(endUrl).String()
	// fmt.Println(endpoint)

	req, err := http.NewRequest(http.MethodGet,
		endpoint,
		// "http://localhost:9090/api/v1/restricted/books?id=12", // これでもダメか...
		// strings.NewReader(values.Encode())
		nil)
	if err != nil {
		t.Error(err)
	}
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	// req.Header.Set("Authorization", "Bearer"+token)
	res := httptest.NewRecorder()
	// res := new(http.ResponseWriter)
	// fmt.Println(c.Get("id"))
	// _, err = http.Get(endpoint)
	// if err != nil {
	// 	t.Error(err)
	// }
	// httptest.NewRequest()
	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(testuserid)
	// err = BookIndex(c) // テーブル分割につき...
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res.Result().StatusCode)
	// fmt.Println(res.Result().Body)
	a := res.Body.Bytes()
	var result []db.Book
	json.Unmarshal(a, &result)
	createdbookid = strconv.Itoa(result[0].Id)
}

func TestBookDelete(t *testing.T) {
	// e := echo.New()
	// base, err := url.Parse("http://localhost:9090")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	endUrl, err := url.Parse("/api/v1/restricted/books/")
	if err != nil {
		fmt.Println(err)
	}
	endpoint := base.ResolveReference(endUrl).String()
	fmt.Println(endpoint)
	req, err := http.NewRequest(http.MethodDelete,
		endpoint,
		nil)
	if err != nil {
		t.Error(err)
	}
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	// req.Header.Set("Authorization", "Bearer"+token)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	// c.SetParamNames()
	// あああああああecho公式に書いてるううううううううああああああ
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(createdbookid)
	// err = BookDelete(c) // テーブル分割につき...// あれ...これって外部からhttpでアクセスしたことになってる...？？？
	// 内部で渡してるだけじゃない？？
	if err != nil {
		t.Error(err)
	}
}

func TestUserDelete(t *testing.T) {
	u := new(db.User)
	u.Name = testname
	err := db.User_record_delete(u)
	if err != nil {
		t.Error(err)
	}
}
