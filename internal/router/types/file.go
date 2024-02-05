package types

type FileRes struct {
	URL     string `json:"url" yaml:"URL" toml:"url"`
	Name    string `json:"name" yaml:"Name" toml:"name"`
	Path    string `json:"path" yaml:"Path" toml:"path"`
	Size    int64  `json:"size" yaml:"Size" toml:"size"`
	Mode    string `json:"mode" yaml:"Mode" toml:"mode"`
	ModTime int64  `json:"mod_time" yaml:"ModTime" toml:"mod_time"`
	IsDir   bool   `json:"is_dir" yaml:"IsDir" toml:"is_dir"`
}

type FileListRes struct {
	Total int        `json:"total" yaml:"Total" toml:"total"`
	Files []*FileRes `json:"files" yaml:"Files" toml:"files"`
}
