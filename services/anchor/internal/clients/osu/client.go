package osu

import (
	"io"

	"github.com/osuTitanic/titanic/services/anchor/internal/clients/bancho"
	"github.com/osuTitanic/titanic/services/anchor/internal/clients/capabilities"
)

type OsuClient struct {
	bancho.ClientBase

	capabilities.Presence
	capabilities.Stats
	capabilities.Friends
	capabilities.Notifications

	capabilities.Chat
	capabilities.Channels

	capabilities.Spectator
	capabilities.SpectatorHost

	capabilities.Multiplayer
	capabilities.MultiplayerObserver

	adapter io.ReadWriteCloser
}

// TODO: implementation
