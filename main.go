package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/zorkian/go-datadog-api"
)

const (
	github string = "GITHUB"
	low string = "LOW"
)

func main() {
	ddApiKey := os.Getenv("INPUT_DD_API_KEY")
	if ddApiKey == "" {
		log.Fatalf("input DD_API_KEY must be set")
	}
	ddAppKey := os.Getenv("INPUT_DD_APP_KEY")
	if ddAppKey == "" {
		log.Fatalf("input DD_APP_KEY must be set")
	}
	client := datadog.NewClient(ddApiKey, ddAppKey)

	eventTitle := os.Getenv("INPUT_EVENT_TITLE")
	if eventTitle == "" {
		log.Fatalf("input EVENT_TITLE must be set")
	}

	eventTagList := strings.Split(
		os.Getenv("INPUT_EVENT_TAGS"),
		",")

	// Event Time in UTC
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	fmt.Println("ZONE : ", loc, " Time : ", now) // UTC

	unixNow := int(now.Unix())
	source := github
	event, err := client.PostEvent(&datadog.Event{
		Id:          nil,  // Optional
		Title:       &eventTitle,
		Text:        nil,  // Optional
		Time:        &unixNow,
		Priority:    nil,  // Optional
		AlertType:   nil,  // Optional
		Host:        nil,  // Optional
		Aggregation: nil,  // Optional
		SourceType:  &source,
		Tags:        eventTagList,  // Optional
		Url:         nil,  // Optional
		Resource:    nil,  // Optional
		EventType:   nil,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	output := fmt.Sprintf("Event: %+v", event)
	fmt.Println(fmt.Sprintf(`::set-output name=output::%s`, output))

}
