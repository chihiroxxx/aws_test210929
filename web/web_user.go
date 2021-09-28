package web

import (
	"backend/db"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
	"unicode/utf8"

	"net/http"

	// "net/url"
	// "strconv"

	"github.com/dgrijalva/jwt-go"
	// "github.com/dgrijalva/jwt-go/request"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate(c echo.Context) error {

	p := c.FormValue("name")
	// p, _ := c.FormParams()
	fmt.Println(p)

	// バリデーション検証（ストロングパラメーター）

	err := userRecordValidation(c) // name 16以内 password 8文字以上
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
		// return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		// これはwebからでもエラーになるの...？？？
		// return c.JSON(http.StatusBadRequest, "don't create...")
	}

	// その前に、早期エラー発見のためにここで、username重複チェックしよう！！

	// バリデーションOKだったらUserRecordCreateする！
	u := newUserRecord(c)
	fmt.Println(&u)
	// if のバリデーション通してOKだったら下の処理
	// で、DBに登録
	resultId, err := db.User_record_create(u) // ここでreturnを加えるか...
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "can't create...")
	}
	if resultId == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, "can't create...") //useridを返す！！
	}
	token, err := getAuth(u, resultId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "can't login...")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":  token, //うーん...これで返せてるのか？？？
		"userid": strconv.Itoa(resultId),
	}) //useridを返す！！
	// return c.JSON(http.StatusOK, resultId) //useridを返す！！
}

// func UserLogout(c echo.Context) error {

// }
func UserLogin(c echo.Context) error {

	u := loginUserRecord(c)
	fmt.Println("まずuserRecordはOK?", u)
	resultId, err := db.User_record_login(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "don't login...")
	}

	token, err := getAuth(u, resultId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "don't login...")
	}
	// fmt.Println(tokenByte)

	// return c.JSON(http.StatusOK, "login OK!!!")
	return c.JSON(http.StatusOK, echo.Map{
		"token":  token, //うーん...これで返せてるのか？？？
		"userid": strconv.Itoa(resultId),
	})
}

// これいつか使いたいなあ...
var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// 違ったわ、この関数使ってないや、、、
func CheckAuth(c echo.Context) error {
	verifyBytes, err := ioutil.ReadFile("./public.key")
	if err != nil {
		fmt.Println(err)
		// return echo.NewHTTPError(http.StatusUnauthorized, "don't authorized and can't login...")
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		fmt.Println(err)
		// return echo.NewHTTPError(http.StatusUnauthorized, "don't authorized and can't login...")
	}

	///ここでrequestの情報をパース（解析するみたい）でチェックするのであろう。
	// request "github.com/dgrijalva/jwt-go/request"  importいる？？？
	// request.ParseFromRequest() // これは *jwt.Token と errorを返す。
	// それをechoに置き換えると...
	token := c.Get("token").(*jwt.Token) // これっぽい！！
	fmt.Println(token)

	return nil
}

/* ここはもうnameとpasswordの一致を確かめた後の
トークン発行だけのところ！！

*/

func getAuth(u *db.User, resultId int) (string, error) { // ひとまず、ただのfuncにしておく！！

	// signBytes, err := ioutil.ReadFile("./secret.key") //ファイル階層注意！！ go.modからの相対パスみたい...
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	// if err != nil {
	// 	fmt.Println(err) // ここでエラーだ！！！ ....つまり...？？？ いや、シークレットキーのファイルが違ったみたいよ...
	// 	/* asn1: structure error: tags don't match
	// 	(16 vs {class:1 tag:15 length:112 isCompound:true})
	// 	{optional:false explicit:false application:false private:false defaultValue:<nil>
	// 		tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pkcs8 @2
	// 	*/
	// }

	// nameとpasswordの一致後！！（念の為もう一回確認してもいいが...）
	// token := jwt.New(jwt.SigningMethodRS256)
	token := jwt.New(jwt.SigningMethodHS256)
	// claim設定！
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = u.Name
	claims["userid"] = resultId
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix() //ひとまず12時間くらいにしてみる...

	// tokenString, err := token.SignedString(signKey)
	tokenString, err := token.SignedString([]byte("secret")) //ここは署名か！！！ あとENVにする...
	// tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY"))) これ！ENVのやつ！これにしたい！
	// あとJWTトークンは どうやら3部構成みたい、 ああ〜この電子署名（署名）があるから、トークンを改ざんできないのか！？！？！
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(tokenString)
	// return []byte(tokenString)

	// このreturnのbyteをechoからのレスポンスに埋め込む...には...
	return tokenString, nil
}

func userRecordValidation(c echo.Context) error {
	// あと正規化もさせたい...小文字英数字のみ と
	u := new(db.User)
	u.Name = c.FormValue("name")
	if utf8.RuneCountInString(u.Name) >= 16 {
		return fmt.Errorf("validation error Name is too long!!! %v", u)
	}
	u.Password_digest = c.FormValue("password_digest")
	if utf8.RuneCountInString(u.Password_digest) < 8 {
		return fmt.Errorf("validation error Password is too short!!! %v", u)
	}
	// ここでnameが重複しているかもチェックさせたい！！！
	// いや、ここで呼ぶと、ぐちゃぐちゃになるからやめる。

	return nil
}
func newUserRecord(c echo.Context) *db.User {
	u := new(db.User)
	// 色々処理------------------
	u.Name = c.FormValue("name")
	//ここでpasswordをハッシュ化しよう！！
	hash, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password_digest")), 12)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hash)
	u.Password_digest = string(hash)
	fmt.Println([]byte(u.Password_digest)) /* stringとしてDBに保存するから、
	その後のログインの時[]byteしたらどうなるのかな？の検証！！
	大丈夫そう！！*/
	// u.Password_digest = c.FormValue("password_digest")

	return u

}

func loginUserRecord(c echo.Context) *db.User {
	u := new(db.User)
	// 色々処理------------------
	u.Name = c.FormValue("name")
	u.Password_digest = c.FormValue("password_digest")

	/*
		bcrypt.CompareHashAndPassword(,[]byte(c.FormValue("password_digest")))

		//ここでpasswordをハッシュ化しよう！！
		hash, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password_digest")), 12)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(hash)
		u.Password_digest = string(hash)
		fmt.Println([]byte(u.Password_digest))
	*/
	/* stringとしてDBに保存するから、
	その後のログインの時[]byteしたらどうなるのかな？の検証！！
	大丈夫そう！！*/
	// u.Password_digest = c.FormValue("password_digest")

	return u
}
