package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	botID string
	forwardSources *ForwardSources
	serverConfig *model.Config
	Utils *Utils
	webhook string
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func (p *Plugin) ForwardSources() *ForwardSources {
	if (p.forwardSources == nil) {
		p.forwardSources = &ForwardSources{}
	}
	return p.forwardSources
}
// See https://developers.mattermost.com/extend/plugins/server/reference/
