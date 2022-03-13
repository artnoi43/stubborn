package cacher

// Config holds default expiration and purge interval in seconds
type Config struct {
	Expiration int `mapstructure:"expiration" json:"expiration"`
	CleanUp    int `mapstructure:"cleanup_interval" json:"cleanupInterval"`
}
