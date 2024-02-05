package types

type Version struct {
	BuildTime string      `json:"build_time" yaml:"BuildTime" toml:"build_time"`
	Version   string      `json:"version" yaml:"Version" toml:"version"`
	Git       VersionGit  `json:"git" yaml:"Git" toml:"git"`
	Go        VersionGO   `json:"go" yaml:"Go" toml:"go"`
	User      VersionUser `json:"user" yaml:"User" toml:"user"`
}

type VersionGit struct {
	Branch string `json:"branch" yaml:"Branch" toml:"branch"`
	Commit string `json:"commit" yaml:"Commit" toml:"commit"`
	URL    string `json:"url" yaml:"URL" toml:"url"`
}

type VersionGO struct {
	Arch    string `json:"arch" yaml:"Arch" toml:"arch"`
	OS      string `json:"os" yaml:"OS" toml:"os"`
	Version string `json:"version" yaml:"Version" toml:"version"`
}

type VersionUser struct {
	Email string `json:"email" yaml:"Email" toml:"email"`
	Name  string `json:"name" yaml:"Name" toml:"name"`
}
