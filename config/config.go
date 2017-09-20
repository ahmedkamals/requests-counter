package config

type (
	Server struct {
		Host, Port, Path string
	}

	Dispatcher struct {
		// Path to the backup file.
		Backup string
	}

	Config struct {
		Server     Server     `json:"server"`
		Dispatcher Dispatcher `json:"dispatcher"`
	}
)

var (
	Configuration *Config
)
