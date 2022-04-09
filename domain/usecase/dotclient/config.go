package dotclient

type Config struct {
	UpStreamTimeout int    `mapstructure:"upstream_timeout" json:"upstreamTimeout"`
	UpStreamIp      string `mapstructure:"upstream_ip" json:"upstreamIp"`
	UpStreamPort    string `mapstructure:"upstream_port" json:"upstreamPort"`
}
