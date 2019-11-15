package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const lennyface = "( ͡° ͜ʖ ͡°)"

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
	// args.Command contains the full command string entered but for now we only do one thing
	return p.API.RegisterCommand(&model.Command{
		Trigger:          "lennyface",
		DisplayName:      "Lennyface",
		Description:      lennyface,
		AutoComplete:     true,
		AutoCompleteDesc: fmt.Sprintf("/lennyface text will appear as \"text %s\"", lennyface),
		AutoCompleteHint: "[text]",
	})
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         fmt.Sprintf("%s %s", strings.TrimPrefix(args.Command, "/lennyface "), lennyface),
	}, nil
}
