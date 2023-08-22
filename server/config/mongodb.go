package config

type MongoDB struct {
	DBname   string   `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	User     string   `mapstructure:"user" json:"user" yaml:"user"`
	Password string   `mapstructure:"password" json:"password" yaml:"password"`
	Hosts    []string `mapstructure:"hosts" json:"hosts" yaml:"hosts"`
}
