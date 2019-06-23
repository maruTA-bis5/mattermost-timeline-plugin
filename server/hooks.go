package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	p.API.LogDebug("MessageHasBeenPosted hook called")
	defer p.API.LogDebug("MessagHasBeenPosted hook exited")

	// detect "goodbye" as "old style" disable command
	message := post.Message
	if message == "@timelime-bot goodbye" {
		p.API.LogDebug("Goodbye command detected")
		p.executeDisableCommand(c, make([]string, 0))
		return
	}

	// check channel is includes forward source
	channelID := post.ChannelId
	channel, _ := p.API.GetChannel(channelID)
	teamID := channel.TeamId
	if !p.ForwardSources().Team(teamID).Contains(channelID) {
		p.API.LogDebug(fmt.Sprintf("Channel(id=%s, name=%s) is not forward sources. return", channelID, channel.Name))
		return
	}
	if channel.Name == "timeline" {
		p.API.LogDebug("Ignore forward loop. return")
		return
	}
	if strings.HasPrefix(post.Type, "system_") {
		p.API.LogDebug("System message. ignore")
		return
	}

	// check #timeline-ignore hashtag
	if strings.Index(message, "#timeline-ignore") != -1 {
		p.API.LogDebug("Forward ignore hashtag found. return")
		return
	}

	timelineChannel, err := p.API.GetChannelByName(teamID, "timeline", false)
	if err != nil {
		p.API.LogError("Timeline channel could not retrieve: " + err.Error())
		return
	}
	originalSender, _ := p.API.GetUser(post.UserId)
	fileIndicator := ""
	if p.Utils.HasFile(post) {
		fileIndicator = " :paperclip:"
	}

	// forward post to timeline channel
	forwardPost := &model.Post{
		UserId:    originalSender.Id, // forward as original sender
		ChannelId: timelineChannel.Id,
		Message:   fmt.Sprintf("%s\r\n(at ~%s)%s", post.Message, channel.Name, fileIndicator),
	}

	if _, err := p.API.CreatePost(forwardPost); err != nil {
		p.API.LogError(err.Error())
	}
}

func (p *Plugin) OnActivate() error {
	p.API.LogDebug("OnActivate hook called")
	defer p.API.LogDebug("OnActivate hook exited")
	if err := p.ensureAccount(); err != nil {
		return err
	}
	p.serverConfig = p.API.GetConfig()

	// TODO load @timeline-bot joined teams & channels
	// TODO ensure timeline channel for teams which @timeline-bot joined

	// TODO Setup slash command: /timeline
	// subcommands
	// - status: check current channel is set as forward source
	// - enable: add current channel to forward sources
	// - disable: remove current channel from forward sources
	// - sources: show forward source channels
	// - debug: show all information about timeline bot

	return nil
}

func (p *Plugin) ensureAccount() error {
	if err := p.migrateIfRequired(); err != nil {
		return err
	}
	botID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "timeline-bot",
		DisplayName: "Timeline Bot",
		Description: "A bot account created by timeline plugin",
	})
	if err != nil {
		return errors.Wrap(err, "Failed to ensure timeline bot.")
	}
	p.botID = botID
	return nil
}

func (p *Plugin) UserHasJoinedTeam(c *plugin.Context, teamMember *model.TeamMember, actor *model.User) {
	p.API.LogDebug("UserHasJoinedTeam hook called")
	defer p.API.LogDebug("UserHasJoinedTeam hook exited")
	if teamMember.UserId != p.botID {
		p.API.LogDebug(fmt.Sprintf("The user(id=%s) is not timeline-bot(id=%s)", teamMember.UserId, p.botID))
		return
	}
	p.ForwardSources().Team(teamMember.TeamId)
}

func (p *Plugin) UserHasLeftTeam(c *plugin.Context, teamMember *model.TeamMember, actor *model.User) {
	p.API.LogDebug("UserHasLeftTeam hook called")
	defer p.API.LogDebug("UserHasLeftTeam hook exited")
	if teamMember.UserId != p.botID {
		p.API.LogDebug(fmt.Sprintf("The user(id=%s) is not timeline-bot(id=%s)", teamMember.UserId, p.botID))
		return
	}
	p.ForwardSources().RemoveTeam(teamMember.TeamId)
}

func (p *Plugin) UserHasJoinedChannel(c *plugin.Context, channelMember *model.ChannelMember, actor *model.User) {
	p.API.LogDebug("UserHasJoinedChannel hook called")
	defer p.API.LogDebug("UserHasJoinedChannel hook exited")
	if channelMember.UserId != p.botID {
		p.API.LogDebug(fmt.Sprintf("The user(id=%s) is not timeline-bot(id=%s)", channelMember.UserId, p.botID))
		return
	}
	channelID := channelMember.ChannelId
	channel, _ := p.API.GetChannel(channelID)
	teamID := channel.TeamId
	p.ForwardSources().Team(teamID).Add(channelID)
}

func (p *Plugin) UserHasLeftChannel(c *plugin.Context, channelMember *model.ChannelMember, actor *model.User) {
	p.API.LogDebug("UserHasLeftChannel hook called")
	defer p.API.LogDebug("UserHasLeftChannel hook exited")
	if channelMember.UserId != p.botID {
		p.API.LogDebug(fmt.Sprintf("The user(id=%s) is not timeline-bot(id=%s)", channelMember.UserId, p.botID))
		return
	}
	channelID := channelMember.ChannelId
	channel, _ := p.API.GetChannel(channelID)
	teamID := channel.TeamId
	p.ForwardSources().Team(teamID).Remove(channelID)
}
