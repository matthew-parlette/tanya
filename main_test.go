package main

import (
	"fmt"
	"testing"

	"github.com/matthew-parlette/houseparty"
)

func TestCreateTodoistTaskFromJiraIssues(t *testing.T) {
	houseparty.ConfigPath = houseparty.GetEnv("CONFIG_PATH", "config")
	houseparty.SecretsPath = houseparty.GetEnv("SECRETS_PATH", "secrets")
	todoistClient, err := houseparty.GetTodoistClient()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	jiraClient, err := houseparty.GetJiraClient()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	count, err := createTodoistTaskFromJiraIssues(todoistClient, jiraClient)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println("Processed", count, "Jira issues")
}

func TestCompleteTodoistTasksFromJiraIssues(t *testing.T) {
	houseparty.ConfigPath = houseparty.GetEnv("CONFIG_PATH", "config")
	houseparty.SecretsPath = houseparty.GetEnv("SECRETS_PATH", "secrets")
	todoistClient, err := houseparty.GetTodoistClient()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	jiraClient, err := houseparty.GetJiraClient()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	count, err := completeTodoistTasksFromJiraIssues(todoistClient, jiraClient)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println("Processed", count, "Jira issues")
}

func TestMain(t *testing.T) {
	t.Skip("Skipping main()")
	// main()
}
