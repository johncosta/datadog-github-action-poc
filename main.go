package main

import (
	"fmt"

	"os"

	"github.com/johncosta/datadog-github-action-poc/datadogext"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	event := datadogext.NewDatadogEvent()
	err := event.Post()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	output := fmt.Sprintf("Event: %+v", event)
	fmt.Printf(`::set-output name=output::%s\n`, output)
}
