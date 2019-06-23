package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	action := ""
	if len(split) > 1 {
		action = split[1]
	}
	params := []string{}
	if len(split) > 2 {
		params = split[2:]
	}
	
	if command != "/timeline" {
		return &model.CommandResponse{}, nil
	}

	// subcommands
	// - status: check current channel is set as forward source
	// - enable: add current channel to forward sources
	// - disable: remove current channel from forward sources
	// - sources: show forward source channels
	// - debug: show all information about timeline bot
	switch action {
	case "status":
		return p.executeStatusCommand(c, params)
	case "enable":
		return p.executeEnableCommand(c, params)
	case "disable":
		return p.executeDisableCommand(c, params)
	case "sources":
		return p.executeSourcesCommand(c, params)
	case "debug":
		return p.executeDebugCommand(c, params)
	}
	return &model.CommandResponse{}, nil
}