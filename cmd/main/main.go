package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

func main() {

	app := &application{}
	app.mattermostConfig = loadConfig()
	app.mattermostClientv4 = model.NewAPIv4Client(app.mattermostConfig.mattermostServer + ":" + app.mattermostConfig.mattermostPort)
	app.mattermostClientv4.SetToken(app.mattermostConfig.mattermostToken)
	//app.mattermostClientv4.SetOAuthToken(app.mattermostConfig.mattermostToken)
	sendMessage(app.mattermostConfig.mattermostChannelID, "Hello, Mattermost!", app)
	app.initWebSocketClient()
	select {}
}

func (app *application) initWebSocketClient() error {
	// Construct the WebSocket URL
	serveru, _ := url.Parse("http://localhost:8065")
	fmt.Println(serveru)
	sir := fmt.Sprintf("ws://%s", "localhost:8065")
	// Create a new WebSocket client
	//sir := fmt.Sprintf("ws://%s", app.config.mattermostServer.Host+app.config.mattermostServer),
	//		app.mattermostClient.AuthToken,
	tempWebSocker, err := model.NewWebSocketClient4(sir, app.mattermostConfig.mattermostToken)
	if err != nil {

		return err
	}
	app.mattermostWebSocketClient = tempWebSocker
	// Start the WebSocket client to begin listening for events
	app.mattermostWebSocketClient.Listen()

	// Create a goroutine to process incoming events
	go func() {
		for event := range app.mattermostWebSocketClient.EventChannel {
			// Handle each WebSocket event
			app.handleWebSocketEvent(event)
		}
	}()
	fmt.Println("WebSocket client initialized")
	return nil
}

func (app *application) handleWebSocketEvent(event *model.WebSocketEvent) {
	// Check if the event is a new posted message
	variable := event.EventType()
	fmt.Println(variable)
	// EventType is posted when a new message is posted in a channel where the bot is listening ( when the bot is added to the channel )
	if event.EventType() == "posted" {
		// Check the message is from Direct Message or Channel
		if event.GetData()["channel_type"] == "D" {
			fmt.Println("Message is from Direct Message")
		} else {
			fmt.Println("Message is from Channel")
			// Check if the message is from the allowed channels
			if contains(app.mattermostConfig.mattermostAllowedChannels, event.GetData()["channel_name"].(string)) {
				message, channelId := parseResponse(event)
				//Check if the message is from the bot
				if strings.HasPrefix(message, app.mattermostConfig.mattermostBotName) {
					fmt.Println(message)
				} else {
					fmt.Println("Message does not start with bot name")
					sendMessage(channelId, "Message does not start with bot name!", app)
				}
			}

		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func parseResponse(event *model.WebSocketEvent) (string, string) {
	//Access the post data
	postJSON, ok := event.GetData()["post"].(string)
	if !ok {
		log.Fatal("parseResponse failed; post key is not a string")
	}

	// Unmarshal the post data into a map
	var postMap map[string]interface{}
	err := json.Unmarshal([]byte(postJSON), &postMap)
	if err != nil {
		log.Fatal(err)
	}
	// Access the message field
	message, ok := postMap["message"].(string)
	if !ok {
		log.Fatal("parseResponse failed; message field is not a string")
	}

	channelid, ok := postMap["channel_id"].(string)
	if !ok {
		log.Fatal("parseResponse failed; channel_id field is not a string")
	}

	return message, channelid
}

func sendMessage(mattermostChannelId string, message string, app *application) {
	post := &model.Post{
		ChannelId: app.mattermostConfig.mattermostChannelID,
		Message:   "Hello, Mattermost!",
	}
	app.mattermostClientv4.CreatePost(post)
	_, response, createPosterror := app.mattermostClientv4.CreatePost(post)
	if createPosterror != nil {
		log.Fatalf("Failed to send message: %v", response.StatusCode)
	}
}
