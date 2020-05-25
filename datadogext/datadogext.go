// package datadogext provides an abstraction of the datadog golang api.  Environment variables provided by github
// actions are loaded, producing an event, which is then sent to Datadog.
package datadogext

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/zorkian/go-datadog-api"
)

const (
	estring string = "error"
	warning string = "warning"
	info    string = "info"
	success string = "success"
	Github  string = "GITHUB"
)

// NewDatadogClient loads the environment variables required for configuration and returns an instantiated DD Client Struct.
func NewDatadogClient() *datadog.Client {
	ddAPIKey := os.Getenv("INPUT_DD_API_KEY")
	if ddAPIKey == "" {
		log.Fatalf("input DD_API_KEY must be set")
	}

	ddAppKey := os.Getenv("INPUT_DD_APP_KEY")
	if ddAppKey == "" {
		log.Fatalf("input DD_APP_KEY must be set")
	}

	return datadog.NewClient(ddAPIKey, ddAppKey)
}

type datadogEvent struct {
	client *datadog.Client
	event  *datadog.Event
}

// GetSource returns the SourceType of the event
func (dde datadogEvent) GetSource() *string {
	return dde.event.SourceType
}

// GetTime returns the time of the event in UnixTime
func (dde datadogEvent) GetTime() *int {
	return dde.event.Time
}

// GetTitle returns the Title of the event
func (dde datadogEvent) GetTitle() *string {
	return dde.event.Title
}

// GetTags returns the Tags on the event
func (dde datadogEvent) GetTags() []string {
	return dde.event.Tags
}

// GetStatus returns the status of the event
func (dde datadogEvent) GetStatus() *string {
	return dde.event.AlertType
}

func (dde *datadogEvent) setSource(source string) {
	dde.event.SourceType = &source
}

func (dde *datadogEvent) setTitle(title string) {
	if title == "" {
		log.Fatalf("input EVENT_TITLE must be set")
	}
	dde.event.Title = &title
}

func (dde *datadogEvent) setTimeToNow() {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	unixNow := int(now.Unix())
	dde.event.Time = &unixNow
}

func (dde *datadogEvent) setTagList(tags string) {
	dde.event.Tags = strings.Split(tags, ",")
}

func (dde *datadogEvent) setStatus(status string) {
	status = strings.ToLower(status)
	switch status {
	case estring, warning, info, success: //valid
		statusString := strings.Join([]string{"status", status}, ":")
		dde.event.AlertType = &statusString
	}
}

// NewDatadogEvent retreives inputs from the environment and returns a constructed event.
func NewDatadogEvent() *datadogEvent {
	client := NewDatadogClient()
	event := &datadogEvent{client, &datadog.Event{}}
	event.setSource(Github)
	event.setTimeToNow()
	event.setTitle(os.Getenv("INPUT_EVENT_TITLE"))
	event.setTagList(os.Getenv("INPUT_EVENT_TAGS"))
	event.setStatus(os.Getenv("INPUT_EVENT_STATUS"))

	return event
}

// Post calls the Datadog api, creating the event.
func (dde datadogEvent) Post() error {
	_, err := dde.client.PostEvent(dde.event)
	if err != nil {
		return err
	}

	return nil
}
