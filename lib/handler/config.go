package handler

type Config struct {
	Address  string `mapstructure:"address" json:"address"`
	Protocol string `mapstructure:"protocol" json:"protocol"`
	// AllTypes means that the server will get DNS records of ALL types for a particular domain name
	AllTypes bool `mapstructure:"all_types" json:"allTypes"`
}
