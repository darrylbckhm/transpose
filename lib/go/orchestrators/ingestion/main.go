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

type server struct{}

// Create struct to hold info about new item
type PropertyItem struct {
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
    desc            string
}

func scanDDB() {
    /*
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
    */
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": }`))
    })
    http.ListenAndServe(":80", nil)
}

func NewDynamoDbRepo() {
    // Initialize a session that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials
    // and region from the shared configuration file ~/.aws/config.
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    // Create DynamoDB client
    svc := dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:4566")})

    // Build the query input parameters
    params := &dynamodb.PutItemInput{
        TableName:                 aws.String("local.test"),
    }

    // Make the DynamoDB Query API call
    result, err := svc.PutItem(params)
    if err != nil {
        fmt.Println("Query API call failed:")
        fmt.Println((err.Error()))
        os.Exit(1)
    }
    fmt.Println(result)
    tableInput := &dynamodb.CreateTableInput{
        AttributeDefinitions: []*dynamodb.AttributeDefinition{
            {
                AttributeName: aws.String("id"),
                AttributeType: aws.String("S"),
            },
        },
        KeySchema: []*dynamodb.KeySchemaElement{
            {
                AttributeName: aws.String("id"),
                KeyType:       aws.String("HASH"),
            },
        },
        ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(10),
            WriteCapacityUnits: aws.Int64(10),
        },
        TableName: aws.String("Teams"),
    }

    _, tableErr := svc.CreateTable(tableInput)

    if tableErr != nil {
        fmt.Println("Got error calling CreateTable:")
        fmt.Println(tableErr.Error())
    }
}

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

        /*
        p := PropertyItem {
                url: _p["url"],
                title: _p["title"],
                image: _p["image"],
                price: _p["price"],
                status: _p["status"],
                listing_type: _p["listing_type"],
                land_listing: _p["land_listing"],
                images: _p["images"],
                property_type: _p["property_type"],
                location: _p["location"],
                location_type: _p["location_type"],
                view_type: _p["view_type"],
                desc: _p["desc"],
        }*/
        fmt.Println(_p)
        res, err := json.Marshal(_p)

        if err != nil {
            fmt.Println(err)
        }

        fmt.Println(string(res))
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
