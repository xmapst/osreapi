package exec

import (
	"time"
)

type Option func(*Cmd)

func WithLogger(logger Ilogger) Option {
	return func(cmd *Cmd) {
		cmd.logger = logger
	}
}

func WithEnv(env []string) Option {
	return func(cmd *Cmd) {
		cmd.env = env
	}
}

func WithShell(shell string) Option {
	return func(cmd *Cmd) {
		cmd.shell = shell
	}
}

func WithScript(content string) Option {
	return func(cmd *Cmd) {
		cmd.content = content
	}
}

func WithWorkspace(workspace string) Option {
	return func(cmd *Cmd) {
		cmd.workspace = workspace
	}
}

func WithScriptDir(scriptDir string) Option {
	return func(cmd *Cmd) {
		cmd.scriptDir = scriptDir
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(cmd *Cmd) {
		cmd.timeout = timeout
	}
}

func WithConsoleCh(consoleCh chan<- string) Option {
	return func(cmd *Cmd) {
		cmd.consoleCh = consoleCh
	}
}
