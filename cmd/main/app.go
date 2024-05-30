package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

type application struct {
	mattermostConfig          config
	mattermostClientv4        *model.Client4
	mattermostWebSocketClient *model.WebSocketClient
	mattermostUser            *model.User
}
