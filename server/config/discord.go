package config

type Discord struct {
	DISCORD_USER_TOKEN string `json:"DISCORD_USER_TOKEN,omitempty"`
	DISCORD_BOT_TOKEN  string `json:"DISCORD_BOT_TOKEN,omitempty"`
	DISCORD_SERVER_ID  string `json:"DISCORD_SERVER_ID,omitempty"`
	DISCORD_CHANNEL_ID string `json:"DISCORD_CHANNEL_ID,omitempty"`
	VERSION            string `json:"VERSION,omitempty"`
	CB_URL             string `json:"CB_URL,omitempty"`
	MSG_PRFIX          string `json:"MSG_PRFIX,omitempty"`
}
