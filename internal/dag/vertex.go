package dag

import (
	"context"
	"time"
)

type VertexFunc = func(ctx context.Context, gName, vName string) error

type Vertex struct {
	ctx   *iContext
	graph *Graph

	cid     int64      // 临时id
	running bool       // 运行中
	fn      VertexFunc // 函数
	adjs    []*Vertex  // 邻接或相邻顶点
	deps    []*Vertex  // 依赖列表
	ndeps   int64      // 顶点的依赖数量
	root    bool       // 根顶点
}

func NewVertex(name string, fn VertexFunc) *Vertex {
	v := &Vertex{
		ctx: &iContext{
			name: name,
		},
		fn: fn,
	}
	return v
}

// Name 获取名称
func (g *Vertex) Name() string {
	return g.ctx.name
}

// Kill 强杀
func (g *Vertex) Kill() error {
	if g.ctx.baseCancel == nil {
		return ErrKill
	}

	g.ctx.baseCancel()

	return nil
}

// Pause 挂起
func (g *Vertex) Pause(duration string) error {
	g.ctx.Lock()
	defer g.ctx.Unlock()

	if g.ctx.controlCtx != nil {
		// 重复挂起, 直接返回
		return nil
	}

	g.ctx.controlCtx, g.ctx.controlCancel = context.WithCancel(context.Background())
	d, err := time.ParseDuration(duration)
	if err == nil && d > 0 {
		g.ctx.controlCtx, g.ctx.controlCancel = context.WithTimeout(context.Background(), d)
	}

	return nil
}

// Resume 解挂
func (g *Vertex) Resume() {
	g.ctx.Lock()
	defer g.ctx.Unlock()
	if g.ctx.controlCancel == nil {
		// 没有挂起不需要恢复,直接返回
		return
	}

	// 解除挂起
	g.ctx.controlCancel()
}

// WaitResume 等待解挂
func (g *Vertex) WaitResume() {
	g.ctx.Lock()
	defer g.ctx.Unlock()

	if g.ctx.controlCtx == nil {
		// 没有挂起不需要d等待,直接返回
		return
	}
	<-g.ctx.controlCtx.Done()
}

// Paused 是否挂起
func (g *Vertex) Paused() bool {
	return g.ctx.controlCtx != nil
}

// WithDeps 为顶点添加依赖顶点。它会检查依赖顶点是否已经在图形中存在，如果不存在，则将依赖顶点添加到图型中
func (v *Vertex) WithDeps(vv ...*Vertex) {
	// adds deps that are not added in graph
	for _, task := range vv {
		if task.cid == 0 {
			v.graph.AddVertex(task)
		}
	}
	v.deps = append(v.deps, vv...)
}
