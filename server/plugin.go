package main

import (
	"fmt"
	"strings"
	"sync"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func (p *Plugin) OnActivate() error {
	// args.Command contains the full command string entered
	return p.API.RegisterCommand(&model.Command{
		Trigger:          "define",
		DisplayName:      "Urban Dictionary",
		Description:      "Urban Dictionary plugin which provides the /define command",
		AutoComplete:     true,
		AutoCompleteDesc: "/define word will look up word on Urban Dictionary and insert its definition into chat.",
	})
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	return UrbanDictionary(args.Command)
}

type UrbanDictionaryResult struct {
	List []UrbanDictionaryEntry `json:"list"`
}

type UrbanDictionaryEntry struct {
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

func UrbanDictionary(in string) (*model.CommandResponse, *model.AppError) {
	trimmed := strings.TrimPrefix(in, "/define ")
	in = url.QueryEscape(trimmed)

	resp, err := http.Get(fmt.Sprintf("http://api.urbandictionary.com/v0/define?term=%s", in))
	if err != nil {
		return nil, model.NewAppError("/define making a request got an error", "", nil, err.Error(), resp.StatusCode)
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, model.NewAppError("/define couldn't read the raw response data from Urban Dictionary", "", nil, err.Error(), resp.StatusCode)
	}

	var result UrbanDictionaryResult
	err = json.Unmarshal(payload, &result)
	if err != nil {
		return nil, model.NewAppError("/define couldn't unmarshal the response JSON from Urban Dictionary", string(payload), nil, err.Error(), resp.StatusCode)
	}

	if len(result.List) == 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Sorry! No results found for \"%s\"", trimmed),
		}, nil
	}

	content := result.List[0]
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         fmt.Sprintf("*****\n**%s**:\n\n%s\n\nExample(s):\n\n%s\n*****\n", trimmed, content.Definition, content.Example),
	}, nil
}
