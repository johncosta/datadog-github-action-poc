package main

import (
	"fmt"
	"log"

	"github.com/johncosta/datadog-github-action-poc/datadogext"
)

func main() {
	event := datadogext.NewDatadogEvent()
	err := event.Post()
	if err != nil {
		log.Fatal(err.Error())
	}

	output := fmt.Sprintf("Event: %+v", event)
	fmt.Println(fmt.Sprintf(`::set-output name=output::%s`, output))
}
