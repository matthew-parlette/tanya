package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	chat "github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
	jira "github.com/andygrunwald/go-jira"
	"github.com/matthew-parlette/houseparty"
	"github.com/sachaos/todoist/lib"
)

func todoistPriority(jiraPriority string) int {
	switch jiraPriority {
	case "Lowest":
		return 2
	case "Low":
		return 2
	case "Medium":
		return 3
	case "High":
		return 4
	case "Highest":
		return 4
	default:
		return 3
	}
}

func syncTodoist(todoistClient *todoist.Client) bool {
	if err := todoistClient.Sync(context.Background()); err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func getContentFromJiraIssue(issue jira.Issue) string {
	return fmt.Sprintf("[[%v] %v - parlette.us](%v)", issue.Key, issue.Fields.Summary, fmt.Sprintf("%v/browse/%v", houseparty.Config("jira-url"), issue.Key))
}

func getTodoistWorkingProjectID(todoistClient *todoist.Client) int {
	// fmt.Println("Retrieving project ID...")
	// This should work to get it, but I'm getting an error that GetIDByName is undefined, so I'll do my own search
	// project := todoistClient.Store.Projects.GetIDByName(houseparty.Config("todoist-project"))
	project := 0
	search := houseparty.Config("todoist-project")
	for _, p := range todoistClient.Store.Projects {
		if p.Name == search {
			project = p.GetID()
		}
	}

	if project == 0 {
		log.Fatal("Could not find project with name ", houseparty.Config("todoist-project"))
	}

	return project
}

func findExistingTodoistTask(todoistClient *todoist.Client, content string) []todoist.Item {
	var existing []todoist.Item
	for _, item := range todoistClient.Store.Items {
		if item.Content == content {
			existing = append(existing, item)
		}
	}
	return existing
}

func createTodoistTaskFromJiraIssues(todoistClient *todoist.Client, jiraClient *jira.Client) (int, error) {
	fmt.Println("Creating tasks from Jira issues...")
	count := 0

	project := getTodoistWorkingProjectID(todoistClient)

	issues, _, err := jiraClient.Issue.Search(
		"assignee = currentUser() AND resolution = Unresolved",
		nil)
	if err != nil {
		return count, err
	}

	for _, issue := range issues {
		// fmt.Println("Processing", issue.Key, "...")
		content := getContentFromJiraIssue(issue)
		existing := findExistingTodoistTask(todoistClient, content)
		if len(existing) > 0 {
			// fmt.Println("Found existing tasks, moving on...")
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
			count = count + 1
		}
	}
	fmt.Println("Created", count, "new tasks from", len(issues), "Jira issues assigned to user")
	return count, nil
}

func completeTodoistTasksFromJiraIssues(todoistClient *todoist.Client, jiraClient *jira.Client) (int, error) {
	fmt.Println("Completing tasks from Jira issues...")

	issues, _, err := jiraClient.Issue.Search(
		"assignee = currentUser() AND resolution != Unresolved AND updated >= startOfMonth(-1)",
		nil)
	if err != nil {
		return 0, err
	}

	items := []int{}
	for _, issue := range issues {
		content := getContentFromJiraIssue(issue)
		existing := findExistingTodoistTask(todoistClient, content)
		for _, task := range existing {
			items = append(items, task.GetID())
		}
	}
	fmt.Println("Found", len(issues), "Jira issues updated since the beginning of last month")
	if len(items) > 0 {
		fmt.Println("Closing", len(items), "todoist tasks...")
		if err = todoistClient.CloseItem(context.Background(), items); err != nil {
			return 0, err
		}
	} else {
		fmt.Println("No todoist tasks found to close, moving on...")
	}
	return len(items), nil
}

func updateOverdueTasks(todoistClient *todoist.Client) (int, error) {
	fmt.Println("Updating overdue tasks to a due date of today...")
	count := 0

	for _, item := range todoistClient.Store.Items {
		if item.DateString != "" && time.Now().UTC().After(item.DateTime()) {
			if strings.Contains(item.DateString, "every") {
				item.DueDateUtc = time.Now().UTC().String()
				item.AllDay = true
			} else {
				item.DateString = "tod"
			}
			if err := todoistClient.UpdateItem(context.Background(), item); err != nil {
				return 0, err
			}
			count = count + 1
		}
	}
	fmt.Println("Found", count, "overdue todoist tasks")
	if count > 0 {
		fmt.Println("Updating due date to today for", count, "todoist tasks...")
	} else {
		fmt.Println("No overdue todoist tasks found, moving on...")
	}
	return count, nil
}

func run(todoistClient *todoist.Client, jiraClient *jira.Client, chatClient *chat.Client) {
	if syncTodoist(todoistClient) {
		created, _ := createTodoistTaskFromJiraIssues(todoistClient, jiraClient)
		completed, _ := completeTodoistTasksFromJiraIssues(todoistClient, jiraClient)
		overdue, _ := updateOverdueTasks(todoistClient)
		if syncTodoist(todoistClient) && (created > 0 || completed > 0 || overdue > 0) {
			message := "I made some changes to your task list\n```"
			if created > 0 {
				message = fmt.Sprintf("%v\n- Created %v tasks from jira issues", message, created)
			}
			if completed > 0 {
				message = fmt.Sprintf("%v\n- Completed %v tasks (jira issue was closed)", message, completed)
			}
			if overdue > 0 {
				message = fmt.Sprintf("%v\n- Updated due date for %v overdue tasks", message, overdue)
			}
			message = fmt.Sprintf("%v\n```", message)
			channel, _ := chatClient.GetChannelId("house-party")
			chatClient.SendMessage(&models.Channel{Id: channel}, message)
		}
	}
	fmt.Printf("Waiting %v seconds to run again...\n", houseparty.Config("interval"))
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
	chatClient, err := houseparty.GetRocketChatClient()
	if err != nil {
		log.Fatal(err)
	}
	houseparty.StartChatListener(chatClient)

	// First run before waiting for ticker
	run(todoistClient, jiraClient, chatClient)

	go func() {
		for {
			select {
			case <-ticker.C:
				run(todoistClient, jiraClient, chatClient)
			case <-shutdown:
				ticker.Stop()
				return
			}
		}
	}()

	// var input string
	// fmt.Scanln(&input)
	// block forever
	<-shutdown
}
