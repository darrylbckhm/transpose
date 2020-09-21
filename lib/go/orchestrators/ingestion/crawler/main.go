package main

import (
    "os"
    "fmt"
    "encoding/json"
    "net/http"
    "strings"
    "time"
	"github.com/gocolly/colly"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)


func main() {
    c := colly.NewCollector(
            colly.AllowedDomains("www.re.cr"),
    )

    c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong:", err)
    })

    c.OnHTML("div.tileItem.visualIEFloatFix", func(e *colly.HTMLElement) {
        _p := make(map[string]string)
        _p["title"] = strings.Trim(e.ChildText("a"), " ")
        _p["url"] = strings.Trim(e.ChildAttr("a", "href"), " ")
        _p["image"] = strings.Trim(e.ChildAttr("img", "src"), " ")
        e.ForEachWithBreak("div", func(_ int, e1 *colly.HTMLElement) bool {
            if e1.Attr("class") == "visualClear" {
                return false
            }
            _key_value := strings.Split(e1.Text, "\n")
            key := strings.ReplaceAll(strings.Trim(strings.ToLower(_key_value[1]), " "), " ", "_")
            _p[key] = strings.Trim(_key_value[2], " ")
            return true
        })
        _p["desc"] = strings.Join(strings.Split(e.ChildText("p"), "\n"), " ")

        res, err := json.Marshal(_p)

        if err != nil {
            fmt.Println(err)
        }
    })

    // Callback for links on scraped pages
	c.OnHTML("div.listingBar", func(e *colly.HTMLElement) {
	    // Extract the linked URL from the anchor tag
	    link := e.ChildAttr("a", "href")
	    // Have our crawler visit the linked URL
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
	})


	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.re.cr/en/costa-rica-real-estate-for-sale/search-properties?form.search-properties.buttons.search=&form.search-properties.widgets.listing_type%3Alist=rs&form.search-properties.widgets.listing_type%3Alist=cs&form.search-properties.widgets.listing_type%3Alist=ll")
}
