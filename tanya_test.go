package tanya

import (
	"fmt"
	"testing"

	"github.com/matthew-parlette/houseparty"
)

func TestTodoistTasksFromJira(t *testing.T) {
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
	count, err := createTodoistTaskfromJiraIssues(todoistClient, jiraClient)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println("Processed", count, "Jira issues")
}

func TestMain(t *testing.T) {
	t.Skip("Skipping main()")
	// main()
}
