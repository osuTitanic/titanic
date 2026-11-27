package irc

import (
	"io"

	"github.com/osuTitanic/titanic/services/anchor/internal/clients/bancho"
	"github.com/osuTitanic/titanic/services/anchor/internal/clients/capabilities"
)

type IrcClient struct {
	bancho.ClientBase

	capabilities.Presence
	capabilities.Notifications

	capabilities.Chat
	capabilities.Channels

	adapter io.ReadWriteCloser
}

// TODO: implementation
