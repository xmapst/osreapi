package starlark

import (
	"context"
	"os"
	"path/filepath"

	"github.com/qri-io/starlib"
	"github.com/segmentio/ksuid"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"

	"github.com/xmapst/osreapi/internal/plugins"
)

func init() {
	plugins.Register(Name, new(Plugin))
}

const Name = "starlark"

type Plugin struct {
}

func (p *Plugin) Name() string {
	return Name
}

func (p *Plugin) Description() string {
	return "starlark plugin provider"
}

func (p *Plugin) WithEnv([]string) plugins.IPlugin {
	return p
}

func (p *Plugin) Run(ctx context.Context, content string) error {
	filename := filepath.Join(os.TempDir(), ksuid.New().String())
	defer os.RemoveAll(filename)
	if err := os.WriteFile(filename, []byte(content), os.ModePerm); err != nil {
		return err
	}
	return p.RunFile(ctx, filename)
}

func (p *Plugin) RunFile(ctx context.Context, filename string) error {
	thread := &starlark.Thread{
		Load: starlib.Loader,
		Name: ksuid.New().String(),
	}
	_, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), thread, filename, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
