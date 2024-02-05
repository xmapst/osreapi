package types

type BaseRes struct {
	Code      int         `json:"code" yaml:"Code" toml:"code" example:"255"`
	Message   string      `json:"message" yaml:"Message" toml:"message" example:"message"`
	Timestamp int64       `json:"timestamp" yaml:"Timestamp" toml:"timestamp"`
	Data      interface{} `json:"data" yaml:"Data" toml:"data"`
}

type Endpoint struct {
	Method string `json:"method" yaml:"Method" toml:"method"`
	Path   string `json:"path" yaml:"Path" toml:"path"`
}
