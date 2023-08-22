package config

type Gpt struct {
	APIKey           string  `json:"api_key"`
	APIURL           string  `json:"api_url"`
	BotDesc          string  `json:"bot_desc"`
	Proxy            string  `json:"proxy"`
	Model            string  `json:"model"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             float32 `json:"top_p"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
	PresencePenalty  float32 `json:"presence_penalty"`
	AuthUser         string  `json:"auth_user"`
	AuthPassword     string  `json:"auth_password"`
}
