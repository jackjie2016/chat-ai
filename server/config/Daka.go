package config

type Daka struct {
	EndHour int      `mapstructure:"end_hour" json:"end_hour" yaml:"end_hour"`
	Groups  []string `mapstructure:"groups" json:"groups" yaml:"groups"`
}
