package main

import (
	"backend/db"
	"backend/scraping"
	"backend/web"
	// "crypto/rsa"
	// "fmt"
	// "io/ioutil"
	"net/http"

	// "os"

	// "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// csv-----------------------!!!!!!!!!!

	db.Db_main()
	// db.User_record_create()
	// file, _ := os.Create("test.txt")
	// file.Write([]byte("testest"))
	// file.Close()
	e := echo.New()
	e.Use(middleware.CORS())
	// user関係
	e.POST("/api/v1/users", web.UserCreate)
	e.POST("/api/v1/login", web.UserLogin)
	// e.POST("/api/v1/logout", web.UserLogout) //単純にReactでtoken削除するか...

	// scraping関係
	e.GET("/api/v1/kino", scraping.Kino)
	e.GET("/api/v1/tsutaya", scraping.Tsutaya)

	// e.POST("/api/v1/login/check", web.CheckAuth)

	// Restricted group ------------test---------------
	// ログイントークンが必要関係！！！
	// r := e.Group("/api/v1/restricted") //名前とかもろもろどうしようかな "/api/v1/login/check"
	r := e.Group("/api/v1/restricted") //名前とかもろもろどうしようかな "/api/v1/login/check"
	// {
	// 	config := middleware.JWTConfig{
	// 		KeyFunc: getKey,
	// 	}

	// r.Use(middleware.JWT(
	// 	func() *rsa.PrivateKey {
	// 	signBytes, err := ioutil.ReadFile("./secret.key") //ファイル階層注意！！ go.modからの相対パスみたい...
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	// 	if err != nil {
	// 		fmt.Println(err) // ここでエラーだ！！！ ....つまり...？？？ いや、シークレットキーのファイルが違ったみたいよ...
	// 		/* asn1: structure error: tags don't match
	// 		(16 vs {class:1 tag:15 length:112 isCompound:true})
	// 		{optional:false explicit:false application:false private:false defaultValue:<nil>
	// 			tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pkcs8 @2
	// 		*/
	// 	}
	// 	return signKey
	// }))

	// r.Use(middleware.JWTWithConfig(config))
	// r.Use(middleware.CORS())
	// r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://localhost:8080"},
	// 	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	// }))

	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted) // <-- ここで判定してるから消さないこと！！

	//ここは...Bookあるなし判定させてからという意味で、、、、
	r.POST("/books", web.BookCreate)

	r.GET("/books/:id", web.BookIndex)
	r.GET("/books/finish/:id", web.BookFinish)

	//現在はこのURLにアクセスはなし！！ コレクションページの時に使う予定！！
	/* // テーブル分割につき...改修中のURLたち
	r.PATCH("/books/:id", web.BookUpdate)
	r.DELETE("/books/:id", web.BookDelete)

	*/
	// thoughtsのcsv！！！
	r.GET("/thoughts/csv/:userid", web.ThoughtCsv)
	// thoughts作成中！
	//そうか、メソッドの違いだけでURLは変えてないんだった！
	r.GET("/thoughts/:id", web.ThoughtIndex)     //接続はOK!!!
	r.DELETE("/thoughts/:id", web.ThoughtDelete) //DELETE!OK!!!
	r.PATCH("/thoughts/:id", web.ThoughtUpdate)  //DELETE!OK!!!

	r.GET("/finishes/:id", web.FinishIndex) //接続はOK!!!

	// 別のweb_total.go に切り分けた！URLもいじる！
	// というか、これいらなそう、reactで計算させたほうがいいっぽいなあ...
	r.GET("/thoughts/total/:id", web.GetMonthTotal)       // totalを計算して必要なとこだけ返してくれるやつ！
	r.GET("/thoughts/total/daily/:id", web.GetDailyTotal) // totalを計算して必要なとこだけ返してくれるやつ！
	// r.GET("/thoughts/total/:id", web.ThoughtGetTotal) // totalを計算して必要なとこだけ返してくれるやつ！

	e.Logger.Fatal(e.Start(":9090"))
}

// ------------------test--------------------------------
/*
func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
*/

func restricted(c echo.Context) error {
	// え... jwtトークン、確認、検証してなくない...？
	// ここで検証させようよ...---------------------------
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// name := claims["name"].(string)
	userid := claims["userid"].(float64)
	// return c.String(http.StatusOK, "Welcome "+name+"!")

	// ここでユーザーIDも返却させたい！！となると
	// claimsから claims["id"].(string) したいから、
	// トークン設定時に、idも含ませないとだな！！！
	return c.JSON(http.StatusOK, int(userid))
}
