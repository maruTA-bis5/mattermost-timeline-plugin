package main

import (
	"github.com/mattermost/mattermost-server/model"
)

func (p *Plugin) migrateIfRequired() error {
	oldUser, err := p.API.GetUserByUsername("timeline-bot")
	if err != nil {
		// no user found -> no migration required
		// TODO handle errors without user not found
		return nil
	}

	if oldUser.IsBot {
		// already bot account -> no migration required
		return nil
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		return err
	}
	for _, team := range teams {
		err = p.migrateForTeam(team, oldUser.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) migrateForTeam(team *model.Team, userID string) *model.AppError {
	channels, err := p.API.GetChannelsForTeamForUser(team.Id, userID, false) // ignore deleted
	if err != nil {
		return err
	}
	for _, channel := range channels {
		p.forwardSources.Team(team.Id).Add(channel.Id)
	}
	return nil
}