package main

import (
	"github.com/deckarep/golang-set"
)

type ForwardSources struct {
	sourcesByTeam map[string]mapset.Set
}

func (fs *ForwardSources) Team(teamID string) mapset.Set {
	fs.ensureMap()
	sources := fs.sourcesByTeam[teamID]
	if sources == nil {
		sources = mapset.NewSet()
		fs.sourcesByTeam[teamID] = sources
	}
	return sources
}

func (fs *ForwardSources) RemoveTeam(teamID string) {
	fs.ensureMap()
	delete(fs.sourcesByTeam,teamID)
}

func (fs *ForwardSources) ensureMap() {
	if fs.sourcesByTeam == nil {
		fs.sourcesByTeam = make(map[string]mapset.Set)
	}
}