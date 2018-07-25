package tanya

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/matthew-parlette/houseparty"
	"github.com/sachaos/todoist/lib"
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
	if err := todoistClient.Sync(context.Background()); err != nil {
		return count, err
	}

	fmt.Println("Retrieving project ID...")
	// This should work to get it, but I'm getting an error that GetIDByName is undefined, so I'll do my own search
	// project := todoistClient.Store.Projects.GetIDByName(houseparty.Config("todoist-project"))
	project := 0
	search := houseparty.Config("todoist-project")
	for _, p := range todoistClient.Store.Projects {
		if p.Name == search {
			project = p.GetID()
		}
	}

	issues, _, err := jiraClient.Issue.Search(
		"assignee = currentUser() AND resolution = Unresolved",
		nil)
	if err != nil {
		return count, err
	}

	for _, issue := range issues {
		fmt.Println("Processing", issue.Key, "...")
		content := fmt.Sprintf("[[%v] %v - parlette.us](%v)", issue.Key, issue.Fields.Summary, fmt.Sprintf("%v/browse/%v", houseparty.Config("jira-url"), issue.Key))
		var existing []todoist.Item
		for _, item := range todoistClient.Store.Items {
			if item.Content == content {
				existing = append(existing, item)
			}
		}
		if len(existing) > 0 {
			fmt.Println("Found existing tasks, moving on...")
			// for _, task := range existing {
			// 	if task.IsOverDueDate() {
			// 		fmt.Println("Task is overdue, updating due date to today...")
			// 		task.DateString = "tod"
			// 	}
			// 	// task.Priority = todoistPriority(issue.Fields.Priority.Name)
			// 	_, err := todoistClient.Item.Update(task)
			// 	if err != nil {
			// 		return count, err
			// 	}
			// }
		} else {
			fmt.Println("Creating new task...")
			item := todoist.Item{}
			item.Content = content
			item.Priority = todoistPriority(issue.Fields.Priority.Name)
			if project > 0 {
				item.ProjectID = project
			}
			item.DateString = "tod"
			if err := todoistClient.AddItem(context.Background(), item); err != nil {
				return count, err
			}
		}
		if err := todoistClient.Sync(context.Background()); err != nil {
			return count, err
		}
		count = count + 1
	}
	return count, nil
}

func main() {
	fmt.Println("Initializing...")
	houseparty.ConfigPath = houseparty.GetEnv("CONFIG_PATH", "config")
	houseparty.SecretsPath = houseparty.GetEnv("SECRETS_PATH", "secrets")
	houseparty.StartHealthCheck()
	interval, err := strconv.Atoi(houseparty.Config("interval"))
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
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
