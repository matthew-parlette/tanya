package main

import (
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/matthew-parlette/houseparty"
)

func TestMain(t *testing.T) {
	houseparty.ConfigPath = houseparty.GetEnv("CONFIG_PATH", "config")
	houseparty.SecretsPath = houseparty.GetEnv("SECRETS_PATH", "secrets")

	todoistClient, err := houseparty.GetTodoistClient()
	if err != nil {
		log.Fatal(err)
	}
	jiraClient, err := houseparty.GetJiraClient()
	if err != nil {
		log.Fatal(err)
	}
	chatClient, err := houseparty.GetRocketChatClient()
	if err != nil {
		log.Fatal(err)
	}

	run(todoistClient, jiraClient, chatClient)
}

func TestTodoistObject(t *testing.T) {
	t.Skip("Skipping Todoist object test")
	houseparty.ConfigPath = houseparty.GetEnv("CONFIG_PATH", "config")
	houseparty.SecretsPath = houseparty.GetEnv("SECRETS_PATH", "secrets")

	todoistClient, err := houseparty.GetTodoistClient()
	if err != nil {
		log.Fatal(err)
	}
	if syncTodoist(todoistClient) {
		spew.Dump(todoistClient.Store.Items[0])
	}
}
