package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

type Utils struct {
	plugin *Plugin
}

func (u *Utils) GetUsername(post *model.Post, sender *model.User) string {
	post.MakeNonNil()
	if u.isWebhookPost(post) {
		overrideName := post.Props["override_username"]
		if overrideName == nil {
			overrideName = "webhook"
		}
		return overrideName.(string)
	}
	return sender.Username
}

func (u *Utils) GetIcon(post *model.Post, sender *model.User) string {
	post.MakeNonNil()
	if u.isWebhookPost(post) {
		overrideIcon := post.Props["override_icon"]
		if overrideIcon == nil {
			// TODO
		}
		return overrideIcon.(string)
	}
	siteURL := u.plugin.serverConfig.ServiceSettings.SiteURL
	return fmt.Sprintf("%s/api/v4/users/%s/image", siteURL, sender.Id)
}

func (u *Utils) isWebhookPost(post *model.Post) bool {
	// TODO
	return false
}

func (u *Utils) HasFile(post *model.Post) bool {
	return len(post.FileIds) >= 1
}