package project_test

import (
	"testing"

	"github.com/apex/apex/project"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetHandler(discard.New())
}

func TestProject_Open_requireConfigValues(t *testing.T) {
	p := &project.Project{
		Path: "../fixtures/project/invalidName",
		Log:  log.Log,
	}
	nameErr := p.Open()

	assert.Contains(t, nameErr.Error(), "Name: zero value")
}
