package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	loggly "github.com/jamespearly/loggly"
)

type Items struct {
	Items []Data
}
type Data struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"keywords"`
	Center      string `json:"center"`
	ID          string `json:"nasa_id"`
}

func main() {

	var tag string
	tag = "My-Go-Demo"

	client := loggly.New(tag)

	err := client.EchoSend("info", "Good morning!")
	fmt.Println("err:", err)

	err = client.Send("error", "Good morning! No echo.")
	fmt.Println("err:", err)

	resp, err := http.Get("https://images-api.nasa.gov/search?q=mercury&year_start=2020&media_type=image")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got the json")

	defer resp.Body.Close()

	byteValue, err := ioutil.ReadAll(resp.Body)

	fmt.Println(byteValue)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Saved the json")

	var items Items

	json.Unmarshal(byteValue, &items)

	fmt.Println(items)

	fmt.Println("Unmarshaled the json")

	fmt.Println("Before for loop")
	for i := 0; i < len(items.Items); i++ {
		fmt.Println("Title: " + items.Items[i].Title)
		fmt.Println("Description: " + items.Items[i].Description)
		fmt.Println("Date: " + items.Items[i].Date)
		fmt.Println("Center: " + items.Items[i].Center)
		fmt.Println("ID: " + items.Items[i].ID)

	}

	fmt.Println("End of File")

}
