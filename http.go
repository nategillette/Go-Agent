package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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
	Title       string `json:"title"`
	Center      string `json:"center"`
	Description string `json:"description_508"`
	ID          string `json:"nasa_id"`
}

type Links struct {
	Link string `json:"href"`
}

func main() {
	duration := time.Duration(1) * time.Second

	tk := time.NewTicker(duration)

	i := 0
	for range tk.C {
		i++
		Query()
		if i > 5 {
			tk.Stop()
			break
		}
	}
}

func Query() {
	var tag string
	tag = "Go"

	client := loggly.New(tag)

	resp, err := http.Get("https://images-api.nasa.gov/search?q=mercury&year_start=2020&media_type=image")
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

		}

	}

	err = client.Send("info", "Succesful Query")
	fmt.Println("err: ", err)
}
