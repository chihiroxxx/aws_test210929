package web

import (
	"backend/db"
	"fmt"
	"net/http"
	"sort"

	"github.com/labstack/echo"
)

// というか、月 週 日 でそれぞれ作ろう！！！

// その週とか月とかの読んだページ数とか集計させる。
//
func GetMonthTotal(c echo.Context) error { //なんにせよユーザIDはいる！！
	id := c.Param("id")                       // これもユーザーIDよ！
	index, err := db.Thought_record_index(id) //結局渡すのはuseridで
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "no record...")
	}
	// んでこのindexを加工してreactに返す
	// んーじゃあcreated_at ->というなの dateから算出しようか
	// type Total struct { //今回欲しいものたち
	// 	date string
	// 	page int
	// }
	// var datearr []Total
	// var jm []byte
	mapin := make(map[string]int)
	if len(index) > 0 { // map(key value)って すんごい便利だなあ〜！！！もっと使いたい！！！map すき！
		for i := 0; i < len(index); i++ {
			//月とか、日とか、結局このスライスで操作してるだけか...
			// mapin[index[i].Date[0:10]] += index[i].Page
			mapin[index[i].Date[0:7]] += index[i].Page

		}
		// fmt.Println(mapin)
		// jsonmap, _ := json.Marshal(mapin) //map -> jsonにマーシャルしなくてよかったのか...
		// fmt.Println(string(jsonmap))
	}
	var arr []Total
	for d, p := range mapin {
		arr = append(arr, Total{Date: d, Page: p})
	}

	fmt.Println(arr)
	return c.JSON(http.StatusOK, arr)
}

func GetDailyTotal(c echo.Context) error { //なんにせよユーザIDはいる！！
	id := c.Param("id")                       // これもユーザーIDよ！
	index, err := db.Thought_record_index(id) //結局渡すのはuseridで
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "no record...")
	}
	// んでこのindexを加工してreactに返す
	// んーじゃあcreated_at ->というなの dateから算出しようか
	// type Total struct { //今回欲しいものたち
	// 	date string
	// 	page int
	// }
	// var datearr []Total
	// var jm []byte
	mapin := make(map[string]int)
	if len(index) > 0 { // map(key value)って すんごい便利だなあ〜！！！もっと使いたい！！！map すき！
		for i := 0; i < len(index); i++ {
			//月とか、日とか、結局このスライスで操作してるだけか...
			// mapin[index[i].Date[0:10]] += index[i].Page
			mapin[index[i].Date[0:10]] += index[i].Page

		}
		// fmt.Println(mapin)
		// jsonmap, _ := json.Marshal(mapin) //map -> jsonにマーシャルしなくてよかったのか...
		// fmt.Println(string(jsonmap))
	}
	var arr []Total
	for d, p := range mapin {
		arr = append(arr, Total{Date: d, Page: p})
	}
	// うまくソートできてないな... どうしようかな...
	sort.Slice(arr, func(i, j int) bool { return index[i].Date < index[j].Date })
	fmt.Println(arr)
	return c.JSON(http.StatusOK, arr)
}

type Total struct {
	Date         string `json:"date"`
	Page         int    `json:"page"`
	Reading_time int    `json:"readingtime"`
}
