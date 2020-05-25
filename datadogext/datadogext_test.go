package datadogext_test

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/johncosta/datadog-github-action-poc/datadogext"
)

func TestDatadogEvent(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	event := datadogext.NewDatadogEvent()

	assert.Equal(t, *event.GetSource(), datadogext.Github)
	assert.NotEmpty(t, event.GetTime())
	assert.Equal(t, *event.GetTitle(), "Test Event Title")
	assert.Equal(t, event.GetTags(), []string{"app:TestApp", "env:Test"})
	assert.Equal(t, *event.GetStatus(), "info")
}
