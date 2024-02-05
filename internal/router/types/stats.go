package types

// State represents activity status of Go.
type State struct {
	Os   StateOs   `json:"os" yaml:"Os" toml:"os"`
	Cpu  StateCpu  `json:"cpu" yaml:"Cpu" toml:"cpu"`
	Ram  StateRam  `json:"ram" yaml:"Ram" toml:"ram"`
	Disk StateDisk `json:"disk" yaml:"Disk" toml:"disk"`
}

type StateOs struct {
	GOOS         string `json:"go_os" yaml:"GOOS" toml:"go_os"`
	NumCPU       int    `json:"num_cpu" yaml:"NumCPU" toml:"num_cpu"`
	Compiler     string `json:"compiler" yaml:"Compiler" toml:"compiler"`
	GoVersion    string `json:"go_version" yaml:"GoVersion" toml:"go_version"`
	NumGoroutine int    `json:"num_goroutine" yaml:"NumGoroutine" toml:"num_goroutine"`
}

type StateCpu struct {
	Cpus  []float64 `json:"cpus" yaml:"Cpus" toml:"cpus"`
	Cores int       `json:"cores" yaml:"Cores" toml:"cores"`
}

type StateRam struct {
	UsedMB      int `json:"used_mb" yaml:"UsedMB" toml:"used_mb"`
	TotalMB     int `json:"total_mb" yaml:"TotalMB" toml:"total_mb"`
	UsedPercent int `json:"used_percent" yaml:"UsedPercent" toml:"used_percent"`
}

type StateDisk struct {
	UsedMB      int `json:"used_mb" yaml:"UsedMB" toml:"used_mb"`
	UsedGB      int `json:"used_gb" yaml:"UsedGB" toml:"used_gb"`
	TotalMB     int `json:"total_mb" yaml:"TotalMB" toml:"total_mb"`
	TotalGB     int `json:"total_gb" yaml:"TotalGB" toml:"total_gb"`
	UsedPercent int `json:"used_percent" yaml:"UsedPercent" toml:"used_percent"`
}
