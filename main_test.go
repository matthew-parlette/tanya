package main

import (
	"log"
	"testing"

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

	if syncTodoist(todoistClient) {
		createTodoistTaskFromJiraIssues(todoistClient, jiraClient)
		completeTodoistTasksFromJiraIssues(todoistClient, jiraClient)
		updateOverdueTasks(todoistClient)
		syncTodoist(todoistClient)
	}
}
