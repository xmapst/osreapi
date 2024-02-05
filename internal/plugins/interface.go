package plugins

import (
	"context"
)

var plugins = make(map[string]IPlugin)

type IPlugin interface {
	Name() string
	Description() string
	WithEnv([]string) IPlugin
	Run(ctx context.Context, content string) error
}

func Register(name string, creator IPlugin) {
	plugins[name] = creator
}

func Get(name string) IPlugin {
	plugin, ok := plugins[name]
	if !ok {
		return nil
	}
	return plugin
}
