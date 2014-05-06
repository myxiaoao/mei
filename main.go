// 倾了全世界的美
package main

import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/PuerkitoBio/gocrawl"
    "github.com/PuerkitoBio/goquery"
    "net/http"
    "regexp"
    "time"
    // "strconv" //这个是为了把int转换为string
)

var rxOk = regexp.MustCompile(`http://douban\.com\/photos\/album\/75978669(\?start(.*))?$`)
// var rxGrep = regexp.MustCompile(`http://www\.douban\.com\/photos\/photo\/\d+$`)

type CustomExtender struct {
    gocrawl.DefaultExtender
}

func first(args ...interface{})interface{} {
    return args[0]
}

func GetConn() *sql.DB {
    db, _ := sql.Open("mysql", "root:root@tcp(localhost:3306)/go_mei?charset=utf8")
    return db
}

func (self *CustomExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
    fmt.Println(ctx.NormalizedURL().String())

    db := GetConn()
    mIns, err := db.Prepare("INSERT INTO mz(photo_href, photo_thumb_src, photo_large_src, photo_public_src, people_href) VALUES( ?, ?, ?, ?, ? )" ) // ? = 占位符

    if err != nil {
        panic(err.Error())
    }

    defer mIns.Close() // main结束是关闭

    //fmt.Println(doc.Find(".photo_wrap").Text())

    doc.Find(".photo_wrap").Each(func(i int, s *goquery.Selection) {
        // For each item found, get the band and title
        // fmt.Println(s.Find("a").First().Attr("title"))
        // fmt.Println(s.Find("a").First().Attr("href"))
        // fmt.Println(s.Find("img").First().Attr("src"))

        var photo_href,photo_thumb_src,photo_large_src,photo_public_src,people_href string
        photo_href = first(s.Find("a").First().Attr("href")).(string)
        photo_thumb_src = first(s.Find("img").First().Attr("src")).(string)
        people_href = first(s.Find("a").First().Attr("title")).(string)

        _, err = mIns.Exec(photo_href, photo_thumb_src, photo_large_src, photo_public_src, people_href)

        // 执行插入
        if err != nil {
            panic(err.Error())
        }
    })

    // if rxGrep.MatchString(ctx.NormalizedURL().String()) {
    // // print problem title
    // fmt.Println(doc.Find("h1").Text())
    // }

    // defer db.Close()

    return nil, true
}

func (self *CustomExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
    // fmt.Println(ctx.NormalizedURL().String())
    return !isVisited && rxOk.MatchString(ctx.NormalizedURL().String())
}

func CustomCrawl() {
    opts := gocrawl.NewOptions(new(CustomExtender))
    opts.CrawlDelay = 3 * time.Second

    c := gocrawl.NewCrawlerWithOptions(opts)
    c.Run("http://www.douban.com/photos/album/75978669/?start=0")
}

func main() {
    CustomCrawl()
}
