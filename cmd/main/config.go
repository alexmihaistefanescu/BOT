package main

type config struct {
	mattermostToken           string
	mattermostServer          string
	mattermostChannelID       string
	mattermostPort            string
	mattermostBotName         string
	mattermostAllowedChannels []string
}

func loadConfig() config {
	var c config
	c.mattermostToken = "zh64nbsne7nytdrcrnfw5a5pqr"                   //os.Getenv("MATTERMOST_TOKEN")
	c.mattermostServer = "http://localhost"                            //os.Getenv("MATTERMOST_SERVER")
	c.mattermostChannelID = "dp9oker9gbdbifxqnmsfw5anya"               //os.Getenv("MATTERMOST_CHANNEL_ID")
	c.mattermostPort = "8065"                                          //os.Getenv("MATTERMOST_PORT")
	c.mattermostBotName = "@test "                                     //os.Getenv("MATTERMOST_BOT_NAME")
	c.mattermostAllowedChannels = []string{"off-topic", "town square"} //strings.Split(os.Getenv("MATTERMOST_ALLOWED_CHANNELS"), ",")
	return c
}
