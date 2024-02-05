package types

type Pool struct {
	Size    int   `json:"size" form:"command_type" yaml:"Size" toml:"size" binding:"required" example:"30"`
	Total   int64 `json:"total" yaml:"Total" toml:"total"`
	Running int64 `json:"running" yaml:"Running" toml:"running"`
	Waiting int64 `json:"waiting" yaml:"Waiting" toml:"waiting"`
}
