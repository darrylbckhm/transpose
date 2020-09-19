package main

import (
    "os"
    "fmt"
    "reflect"
    "net/http"
    "strings"
    "time"
	"github.com/gocolly/colly"
	//"github.com/PuerkitoBio/goquery"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

type server struct{}

// Create struct to hold info about new item
type Item struct {
    url             string
    title           string
    image           string
    price           string
    status          string
    listing_type    string
    land_listing    string
    images          string
    property_type   string
    location        string
    location_type   string
    view_type       string
}

func scanDDB() {
    // Initialize a session that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials
    // and region from the shared configuration file ~/.aws/config.
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Create DynamoDB client
    svc := dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:4566")})

    // Build the query input parameters
    params := &dynamodb.ScanInput{
        TableName:                 aws.String("local.test"),
    }

    // Make the DynamoDB Query API call
    result, err := svc.Scan(params)
    if err != nil {
        fmt.Println("Query API call failed:")
        fmt.Println((err.Error()))
        os.Exit(1)
    }
    fmt.Println(result)

    for _, i := range result.Items {
        fmt.Println(i)
        fmt.Println(i["title"])
        keys := reflect.ValueOf(i).MapKeys()
        fmt.Println(keys)
        item := Item{
            //title: i,
        }
        fmt.Println(item)
    }
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": }`))
    })
    http.ListenAndServe(":80", nil)
}

func main() {
    c := colly.NewCollector(
            colly.AllowedDomains("www.re.cr"),
    )

    c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong:", err)
    })

    c.OnHTML("div.tileItem.visualIEFloatFix", func(e *colly.HTMLElement) {
        title := e.ChildText("a")
        url := e.ChildAttr("a", "href")
        image := e.ChildAttr("img", "src")
        fmt.Println("Title: ", strings.Trim(title, " "))
        fmt.Println("URL: ", strings.Trim(url, " "))
        fmt.Println("Image: ", strings.Trim(image, " "))
        e.ForEachWithBreak("div", func(_ int, e1 *colly.HTMLElement) bool {
            if e1.Attr("class") == "visualClear" {
                return false
            }
            _key_value := strings.Split(e1.Text, "\n")
            var key_value []string
            for elem := 0; elem < len(_key_value); elem++ {
                if strings.Trim(_key_value[elem], " ") != "" {
                    key_value = append(key_value, _key_value[elem])
                }
            }
            if len(key_value) == 1 {
                if key_value[0] == "Property Type" {
                    key_value = append(key_value, "Other")
                }
            }
            fmt.Println("Key: ", strings.Trim(key_value[0], " "))
            fmt.Println("Value: ", strings.Trim(key_value[1], " "))
            return true
        })
        desc := e.ChildText("p")
        fmt.Println("Desc: ", strings.Split(desc, "\n"))
        fmt.Println("")
    })

    // Callback for links on scraped pages
	c.OnHTML("div.listingbar", func(e *colly.HTMLElement) {
        //next := map[string]bool{}
        //e.ForEachWithBreak("span.next", func(_ int, e1 *colly.HTMLElement) bool {
        //    res := strings.Trim(e1, " ")
        //    fmt.Println(res)
        //    next[res] = true
        //}
	    // Extract the linked URL from the anchor tag
	    link := e.ChildAttr("span", "class")
        fmt.Println(link)
	// Have our crawler visit the linked URL
	//	c.Visit(e.Request.AbsoluteURL(link))
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
