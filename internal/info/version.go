package info

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	Version   string
	GitUrl    string
	GitBranch string
	GitCommit string
	BuildTime string
	UserName  string
	UserEmail string
)

func VersionInfo() string {
	var info = []string{
		fmt.Sprintf("Go Version: %s", runtime.Version()),
		fmt.Sprintf("OS: %s", runtime.GOOS),
		fmt.Sprintf("Arch: %s", runtime.GOARCH),
	}
	if Version != "" {
		info = append(info, fmt.Sprintf("Version: %s", Version))
	}
	if GitUrl != "" {
		info = append(info, fmt.Sprintf("Git URL: %s", GitUrl))
	}
	if GitBranch != "" {
		info = append(info, fmt.Sprintf("Git Branch: %s", GitBranch))
	}
	if GitCommit != "" {
		info = append(info, fmt.Sprintf("Git Commit: %s", GitCommit))
	}
	if BuildTime != "" {
		info = append(info, fmt.Sprintf("Build Time: %s", BuildTime))
	}
	if UserName != "" {
		info = append(info, fmt.Sprintf("Developer: %s", UserName))
	}
	if UserEmail != "" {
		info = append(info, fmt.Sprintf("Mail: %s", UserEmail))
	}
	return strings.Join(info, "\n")
}

func PrintHeadInfo() {
	fmt.Println(VersionInfo())
}
