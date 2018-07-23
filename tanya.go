package tanya

import (
	"context"
	"fmt"
	"log"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/matthew-parlette/houseparty"
)

func todoistPriority(jiraPriority string) int {
	switch jiraPriority {
	case "Lowest":
		return 3
	case "Low":
		return 3
	case "Medium":
		return 2
	case "High":
		return 1
	case "Highest":
		return 1
	default:
		return 2
	}
}

func createTodoistTaskfromJiraIssues(todoistClient *todoist.Client, jiraClient *jira.Client) (int, error) {
	count := 0
	ctx := context.Background()
	if err := todoistClient.FullSync(ctx, []todoist.Command{}); err != nil {
		return count, err
	}

	fmt.Println("Retrieving projects...")
	projects := todoistClient.Project.FindByName(houseparty.Config("todoist-project"))

	issues, _, err := jiraClient.Issue.Search(
		"assignee = currentUser() AND resolution = Unresolved",
		nil)
	if err != nil {
		return count, err
	}

	for _, issue := range issues {
		fmt.Println("Processing", issue.Key, "...")
		content := fmt.Sprintf("[[%v] %v - parlette.us](%v)", issue.Key, issue.Fields.Summary, fmt.Sprintf("%v/browse/%v", houseparty.Config("jira-url"), issue.Key))
		existing := todoistClient.Item.FindByContent(content)
		if len(existing) > 0 {
			fmt.Println("Found existing tasks, updating...")
			for _, task := range existing {
				if task.IsOverDueDate() {
					fmt.Println("Task is overdue, updating due date to today...")
					task.DateString = "tod"
				}
				// task.Priority = todoistPriority(issue.Fields.Priority.Name)
				_, err := todoistClient.Item.Update(task)
				if err != nil {
					return count, err
				}
			}
		} else {
			fmt.Println("Creating new task...")
			todoistClient.Item.Add(todoist.Item{
				ProjectID:  projects[0].ID,
				Content:    content,
				DateString: "tod",
				Priority:   todoistPriority(issue.Fields.Priority.Name),
			})
		}
		if err = todoistClient.Commit(ctx); err != nil {
			return count, err
		}
		if err = todoistClient.FullSync(ctx, []todoist.Command{}); err != nil {
			return count, err
		}
		count = count + 1
	}
	return count, nil
}

func main() {
	fmt.Println("Initializing...")
	houseparty.ConfigPath = houseparty.GetEnv("CONFIG_PATH", "/home/matt/src/tanya/config")
	houseparty.SecretsPath = houseparty.GetEnv("SECRETS_PATH", "/home/matt/src/tanya/secrets")
	ticker := time.NewTicker(30 * time.Second)
	shutdown := make(chan struct{})

	todoistClient, err := houseparty.GetTodoistClient()
	if err != nil {
		log.Fatal(err)
	}
	jiraClient, err := houseparty.GetJiraClient()
	if err != nil {
		log.Fatal(err)
	}
	_, err = houseparty.GetRocketChatClient()
	if err != nil {
		log.Fatal(err)
	}

	// First run before waiting for ticker
	_, err = createTodoistTaskfromJiraIssues(todoistClient, jiraClient)
	fmt.Println("Press enter to shutdown...")

	go func() {
		for {
			select {
			case <-ticker.C:
				createTodoistTaskfromJiraIssues(todoistClient, jiraClient)
				fmt.Println("Press enter to shutdown...")
			case <-shutdown:
				ticker.Stop()
				return
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}