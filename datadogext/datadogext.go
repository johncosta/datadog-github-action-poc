// package datadogext provides an abstraction of the datadog golang api.  Environment variables provided by github
// actions are loaded, producing an event, which is then sent to Datadog.
package datadogext

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/levenlabs/golib/timeutil"
	log "github.com/sirupsen/logrus"
	"github.com/zorkian/go-datadog-api"
)

const (
	estring string  = "error"
	warning string  = "warning"
	info    string  = "info"
	success string  = "success"
	Github  string  = "GITHUB"
	count   string  = "count"
	one     float64 = 1
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
	client            *datadog.Client
	event             *datadog.Event
	eventMetric       *float64
	eventMetricName   *string
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
		dde.event.AlertType = &status
	}
}

func (dde *datadogEvent) setEventMetric(status string) {
	if status == "" {
		return
	}

	metric, err := strconv.ParseFloat(status, 64)
	if err != nil {
		log.Fatalf("unable to convert `%s` to int", status)
	}

	dde.eventMetric = &metric
}

func (dde *datadogEvent) setEventMetricName(name string) {
	if name != "" {
		dde.eventMetricName = &name
	}
}

// NewDatadogEvent retreives inputs from the environment and returns a constructed event.
func NewDatadogEvent() *datadogEvent {
	client := NewDatadogClient()
	event := &datadogEvent{client, &datadog.Event{}, nil, nil}
	event.setSource(Github)
	event.setTimeToNow()
	event.setTitle(os.Getenv("INPUT_EVENT_TITLE"))
	event.setTagList(os.Getenv("INPUT_EVENT_TAGS"))
	event.setStatus(os.Getenv("INPUT_EVENT_STATUS"))
	event.setEventMetric(os.Getenv("INPUT_EVENT_METRIC"))
	event.setEventMetricName(os.Getenv("INPUT_EVENT_METRIC_NAME"))
	log.Debugf("New Event Generated and Configured: `%+v`", event)

	return event
}

// Post calls the Datadog api, creating the event.
func (dde datadogEvent) Post() (err error) {
	_, err = dde.client.PostEvent(dde.event)
	if err != nil {
		return err
	}

	log.Info("Event posted successfully.")
	countType := count
	int64Time := int64(*dde.event.Time)
	convertedTime := timeutil.TimestampFromInt64(int64Time).Float64()

	if dde.eventMetric == nil {
		return nil
	}

	sOne := one
	statusMetricName := strings.Join([]string{*dde.eventMetricName, *dde.event.AlertType}, ".")
	err = dde.client.PostMetrics(
		[]datadog.Metric{
			{
				Metric: dde.eventMetricName,
				Tags:   dde.event.Tags,
				Type:   &countType,
				Points: []datadog.DataPoint{
					{
						&convertedTime,
						dde.eventMetric,
					},
				},
			},
			{
				Metric: &statusMetricName,
				Tags:   dde.event.Tags,
				Type:   &countType,
				Points: []datadog.DataPoint{
					{
						&convertedTime,
						&sOne,
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	log.Info("Metric posted successfully.")
	return nil
}
