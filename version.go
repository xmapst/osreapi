package osreapi

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

var (
	Version      string
	GoVersion    string
	GitUrl       string
	GitBranch    string
	GitCommit    string
	GitLatestTag string
	BuildTime    string
)

func VersionIfo() string {
	return fmt.Sprintf("Version: %s\nGoVersion: %s\nGitUrl: %s\nGitBranch: %s\nGitCommit: %s\nGitLatestTag: %s\nBuildTime: %s",
		Version, GoVersion, GitUrl, GitBranch, GitCommit, GitLatestTag, BuildTime)
}

func PrintHeadInfo() {
	figure.NewFigure("OsRE API", "doom", true).Print()
	fmt.Println()
	fmt.Println(VersionIfo())
}
