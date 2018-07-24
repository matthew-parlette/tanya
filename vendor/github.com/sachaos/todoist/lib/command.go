package todoist

import (
	"encoding/json"
	"net/url"

	"github.com/satori/go.uuid"
)

type Command struct {
	Args   interface{} `json:"args"`
	TempID string      `json:"temp_id"`
	Type   string      `json:"type"`
	UUID   string      `json:"uuid"`
}

type Commands []Command

func NewCommand(command_type string, command_args interface{}) Command {
	return Command{
		UUID:   uuid.NewV4().String(),
		TempID: uuid.NewV4().String(),
		Type:   command_type,
		Args:   command_args,
	}
}

func (commands Commands) UrlValues() url.Values {
	commands_text, err := json.Marshal(commands)
	if err != nil {
		return url.Values{}
	}
	return url.Values{
		"commands": {string(commands_text)},
	}
}
