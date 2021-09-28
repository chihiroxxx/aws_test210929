package scraping

import (
	"fmt"
	// "io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
)

type Scraping struct {
	// Id     int64  `json:"id"` //とりあえず削除
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"itemUrl"`
	// Text   string `json:"text"`
	Image string `json:"imageUrl"`
}

func Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello go!!!")
}
func Kino(c echo.Context) error {
	// まずReactからURLparamsで送ることにする
	// https://www.kinokuniya.co.jp/f/dsd---?p=1 &qs=true&ptk=01&  q=react
	q := c.QueryParam("q")
	pg := c.QueryParam("page")
	// 案------------------------------------------------------------------------------------------------------------
	// https://tabelog.com/tokyo/A1303/A130301/rstLst/2/   今のurlね
	// https://tabelog.com/tokyo/A1303/A130301 までReactからparamsで送れば 検索要項もかえられるし、可変性が高そう（適応性？？）
	// var i int // んでforで回せそう
	// i := 2
	// fmt.Println(strconv.Itoa(i))
	// 案------------------------------------------------------------------------------------------------------------
	var arr []Scraping
	// in := 1
	count := 0
	for i := 1; i <= 2; i++ {
		// url := "https://www.kinokuniya.co.jp/f/dsd---?p=" + strconv.Itoa(i) + "&qs=true&ptk=01&q=" + q
		// https://www.kinokuniya.co.jp/disp/CSfDispListPage_001.jsp?qs=true&ptk=01&q=%E3%83%9D%E3%82%B1%E3%83%A2%E3%83%B3&p=2
		// url := "https://www.kinokuniya.co.jp/f/dsd---?p=" + pg + "&qs=true&ptk=01&q=" + q
		pg, _ := strconv.Atoi(pg)
		newpg := pg + count
		url := "https://www.kinokuniya.co.jp/disp/CSfDispListPage_001.jsp?qs=true&ptk=01&q=" + q + "&p=" + strconv.Itoa(newpg)
		// fmt.Println(url)
		doc, err := goquery.NewDocument(url)
		if err != nil {
			fmt.Println(err)
		}
		// doc.Find("h3.heightLine-2").Each(func(index int, s *goquery.Selection) {
		doc.Find("div.list_area").Each(func(index int, s *goquery.Selection) {
			// doc.Find("h3.list-rst__rst-name > a").Each(func(index int, s *goquery.Selection) {
			// str := s.Text()
			// fmt.Println(str)
			var data Scraping
			// data.Id = int64(in) //とりあえず削除
			// str := s.Children().Children().Text()
			str := s.Find("h3.heightLine-2").Text()
			data.Title = str
			if s.Find("h3.heightLine-2 > a").Text() != "" {
				str := s.Find("h3.heightLine-2 > a").Text()
				data.Title = str
			}
			// str := s.Children().Text()
			aut := s.Find("p.clearfix").Text()
			aut = strings.TrimSpace(aut)
			// str := s.Children().Text()
			data.Author = aut
			href, _ := s.Find("h3.heightLine-2 > a").Attr("href")
			data.URL = href
			// dsc := s.Children().Next().Text()
			// data.Text = dsc
			img, err := s.Find("div.listphoto > a.thumbnail_box > img").Attr("src")
			if !err {
				img = "testest"
			} else {
				img = "https://www.kinokuniya.co.jp" + img[2:]

			}
			data.Image = img
			arr = append(arr, data)
			// in++
		})
		count++

	}
	fmt.Println(arr)
	return c.JSON(http.StatusOK, echo.Map{
		"items": arr})
}

func Tsutaya(c echo.Context) error {
	// まずReactからURLparamsで送ることにする
	// https://www.kinokuniya.co.jp/f/dsd---?p=1 &qs=true&ptk=01&  q=react
	// https://tsutaya.tsite.jp/search?dm=0&ds=1&st=0&p= 2&ic=3&  k=react
	q := c.QueryParam("q")
	pg := c.QueryParam("page")
	// 案------------------------------------------------------------------------------------------------------------
	// https://tabelog.com/tokyo/A1303/A130301/rstLst/2/   今のurlね
	// https://tabelog.com/tokyo/A1303/A130301 までReactからparamsで送れば 検索要項もかえられるし、可変性が高そう（適応性？？）
	// var i int // んでforで回せそう
	// i := 2
	// fmt.Println(strconv.Itoa(i))
	// 案------------------------------------------------------------------------------------------------------------
	var arr []Scraping
	// in := 1
	// for i := 1; i <= 10; i++ {
	// url := "https://tsutaya.tsite.jp/search?dm=0&ds=1&st=0&p=" + strconv.Itoa(i) + "&ic=3&k=" + q
	url := "https://tsutaya.tsite.jp/search?dm=0&ds=1&st=0&p=" + pg + "&ic=3&k=" + q
	// fmt.Println(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
	}
	doc.Find("div.c_unit_col-main_in > ul.c_thumb_list_row > li").Each(func(index int, s *goquery.Selection) {
		// doc.Find("li > div > div.c_thumb_info").Each(func(index int, s *goquery.Selection) {
		// doc.Find("h3.list-rst__rst-name > a").Each(func(index int, s *goquery.Selection) {
		// str := s.Text()
		// fmt.Println(str)
		var data Scraping
		// data.Id = int64(in) //とりあえず削除
		// str := s.Children().Children().Text()
		str := s.Find("div > div.c_thumb_info > p.c_thumb_tit").Children().Text()
		// str := s.Children().Text()
		data.Title = str
		aut := s.Find("div > div.c_thumb_info > p.c_thumb_author").Children().Text()
		// aut = strings.TrimSpace(aut)
		// str := s.Children().Text()
		data.Author = aut

		href, _ := s.Find("a").Attr("href")
		data.URL = "https://tsutaya.tsite.jp/" + href
		// dsc := s.Children().Next().Text()
		// data.Text = dsc
		img, _ := s.Find("a > div > span > img").Attr("src")
		// img = "https://www.kinokuniya.co.jp" + img[2:]
		data.Image = img
		arr = append(arr, data)
		// in++
	})

	// }
	fmt.Println(arr)
	// return c.JSON(http.StatusOK, arr)
	// ここswift用に修正する！！！
	// ああ〜やっぱりarrのkeyに"item"つまりkeyで取得すれば取れるわ...
	// これ、react変わっちゃうかな...
	return c.JSON(http.StatusOK, echo.Map{
		"items": arr})
}

// type Page struct {
// 	Title string
// 	Body  []byte
// }

// func (p *Page) save() error {
// 	filename := p.Title + ".txt"
// 	return ioutil.WriteFile(filename, p.Body, 0600)
// }
