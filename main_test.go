package main

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/matthew-parlette/houseparty"
)

func TestRun(t *testing.T) {
	run()
}

func TestChatListener(t *testing.T) {
	houseparty.StartChatListener()
}

func TestTodoistObject(t *testing.T) {
	t.Skip("Skipping Todoist object test")
	if syncTodoist() {
		spew.Dump(houseparty.TodoistClient.Store.Items[0])
	}
}
