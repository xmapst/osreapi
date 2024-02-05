package lua

import (
	"context"
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/plugins"
	"github.com/xmapst/osreapi/internal/utils"
)

func init() {
	plugins.Register(Name, new(Plugin))
}

const Name = "lua"

type Plugin struct {
	env []string
}

func (p *Plugin) Name() string {
	return Name
}

func (p *Plugin) Description() string {
	return "lua plugin provider"
}

func (p *Plugin) WithEnv([]string) plugins.IPlugin {
	return p
}

func (p *Plugin) Run(ctx context.Context, content string) error {
	L := p.newLuaState(nil)
	defer L.Close()

	L.SetContext(ctx)

	// 添加、删除lua脚本参数
	L.SetGlobal("Params", luar.New(L, utils.SliceToStrMap(p.env)))
	// setPrintf
	var outline string
	p.setPrintf(L, &outline)

	beginTime := time.Now()
	logx.Infof("[LuaScropt]running lua script start")
	logx.Debugf("[LuaScropt]lua script content:\n%v", content)
	defer func() {
		logx.Infof("[LuaScropt]running lua script end, cost=%vus", time.Now().Sub(beginTime).Microseconds())
	}()

	if err := L.DoString(content); err != nil {
		logx.Errorf("[LuaScropt]lua script error:%v", err)
		return err
	}
	logx.Infoln(outline)
	return nil
}

func (p *Plugin) newLuaState(globals map[string]interface{}) *lua.LState {
	L := lua.NewState(lua.Options{IncludeGoStackTrace: true})

	// 默认通用Global符号
	L.SetGlobal("Debugf", luar.New(L, logx.Debugf))
	L.SetGlobal("Debug", luar.New(L, logx.Debugln))
	L.SetGlobal("Infof", luar.New(L, logx.Infof))
	L.SetGlobal("Info", luar.New(L, logx.Infoln))
	L.SetGlobal("Warnf", luar.New(L, logx.Warnf))
	L.SetGlobal("Warn", luar.New(L, logx.Warnln))
	L.SetGlobal("Errorf", luar.New(L, logx.Errorf))
	L.SetGlobal("Error", luar.New(L, logx.Errorln))

	for k, v := range globals {
		L.SetGlobal(k, luar.New(L, v))
	}

	return L
}

func (p *Plugin) setPrintf(l *lua.LState, outline *string) {
	printf := func(format string, args ...interface{}) {
		*outline += fmt.Sprintf(format, args...)
	}
	_print := func(args ...interface{}) {
		*outline += fmt.Sprint(args...)
	}
	_println := func(args ...interface{}) {
		*outline += fmt.Sprintln(args...)
	}
	l.SetGlobal("Printf", luar.New(l, printf))
	l.SetGlobal("Print", luar.New(l, _print))
	l.SetGlobal("Println", luar.New(l, _println))
}
