package main

import (
	"fmt"
	"log"
	"os"

	"github.com/johncosta/datadog-github-action-poc/datadogext"
)

func main() {
	event := datadogext.NewDatadogEvent()
	err := event.Post()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	output := fmt.Sprintf("Event: %+v", event)
	fmt.Println(fmt.Sprintf(`::set-output name=output::%s`, output))
}
