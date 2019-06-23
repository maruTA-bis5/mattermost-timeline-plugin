package main

import (
	"github.com/mattermost/mattermost-server/plugin"
)

func main() {
	p := &Plugin{}
	p.Utils = &Utils{
		plugin: p,
	}
	plugin.ClientMain(p)
}
