package types

type Healthyz struct {
	Server string `json:"server" yaml:"Server" toml:"server"`
	Client string `json:"client" yaml:"Client" toml:"client"`
	State  string `json:"state" yaml:"State" toml:"state"`
}
