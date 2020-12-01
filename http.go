package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	loggly "github.com/jamespearly/loggly"
)

type Request struct {
	Collection Collection `json:"collection"`
}

type Collection struct {
	Items []Items `json:"items`
}

type Items struct {
	Data  []Data  `json:"data"`
	Links []Links `json:"links"`
}

type Data struct {
	ID          string `json:"nasa_id"`
	Title       string `json:"title"`
	Center      string `json:"center"`
	Description string `json:"description_508"`
}

type Links struct {
	Link string `json:"href"`
}
type Item struct {
	ID          string `json:"ID"`
	Title       string
	Center      string
	Description string
	URL         string
}

func main() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	timer, _ := strconv.Atoi(os.Args[1])

	if timer <= 0 {
		timer = 60
	}

	duration := time.Duration(timer) * time.Second

	tk := time.NewTicker(duration)

	for range tk.C {
		Query(svc)
	}
}

func Query(svc *dynamodb.DynamoDB) {
	var tag string
	tag = "GoQuery"

	client := loggly.New(tag)

	resp, err := http.Get("https://images-api.nasa.gov/search?q=moon&year_start=2020")
	err = client.Send("info", "Request made")
	if err != nil {
		client.Send("error", "Request failed")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	byteValue, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var request Request

	e := json.Unmarshal(byteValue, &request)

	if e != nil {
		log.Fatal(e)
	}

	for i := 0; i < len(request.Collection.Items); i++ {
		for j := 0; j < len(request.Collection.Items[i].Data); j++ {
			fmt.Println("Title: " + request.Collection.Items[i].Data[j].Title)
			fmt.Println("Center: " + request.Collection.Items[i].Data[j].Center)
			fmt.Println("Description: " + request.Collection.Items[i].Data[j].Description)
			fmt.Println("ID: " + request.Collection.Items[i].Data[j].ID)
			for k := 0; k < len(request.Collection.Items[i].Links); k++ {
				fmt.Println("Link: " + request.Collection.Items[i].Links[k].Link)
			}

			links := ""
			if len(request.Collection.Items[i].Links) > 0 {
				links = request.Collection.Items[i].Links[0].Link
			}

			item := Item{
				ID:          request.Collection.Items[i].Data[j].ID,
				Title:       request.Collection.Items[i].Data[j].Title,
				Center:      request.Collection.Items[i].Data[j].Center,
				Description: request.Collection.Items[i].Data[j].Description,
				URL:         links,
			}
			AddDBItem(item, svc)

		}

	}

	err = client.Send("info", "Succesful Query")
	fmt.Println("err: ", err)
}

func AddDBItem(item Item, svc *dynamodb.DynamoDB) {

	var tag string
	tag = "DynamoDB_Write"

	client := loggly.New(tag)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Error in Marshalling Map")
		client.Send("error", "Couldn't Marshal Map")
		os.Exit(1)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("ngillet2-NASAPhotos"),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err)
		client.Send("error", "Couldn't PutItem")
		os.Exit(1)

	}

	length := len(item.Center) + len(item.Description) + len(item.ID) + len(item.Title) + len(item.URL)

	err = client.Send("info", "New Entry Length:"+strconv.Itoa(length))
	fmt.Println("err: ", err)

}
