package dag

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Graph struct {
	wg  *sync.WaitGroup
	ctx *iContext

	vertex []*Vertex // 顶点集
}

func New(name string) *Graph {
	g := &Graph{
		wg: new(sync.WaitGroup),
		ctx: &iContext{
			name: name,
		},
	}

	// 基础上下文
	g.ctx.baseCtx, g.ctx.baseCancel = context.WithCancel(context.Background())

	// 加入管理
	join(fmt.Sprintf(graphPrefix, g.Name()), g)

	return g
}

// Name 获取名称
func (g *Graph) Name() string {
	return g.ctx.name
}

// Kill 强杀
func (g *Graph) Kill() error {
	if g.ctx.baseCancel == nil {
		return ErrKill
	}

	g.ctx.baseCancel()

	return nil
}

// Pause 挂起
func (g *Graph) Pause(duration string) error {
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
func (g *Graph) Resume() {
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
func (g *Graph) WaitResume() {
	g.ctx.Lock()
	defer g.ctx.Unlock()

	if g.ctx.controlCtx == nil {
		// 没有挂起不需要d等待,直接返回
		return
	}
	<-g.ctx.controlCtx.Done()
}

// Paused 是否挂起
func (g *Graph) Paused() bool {
	return g.ctx.controlCtx != nil
}

// AddVertex 添加顶点
func (g *Graph) AddVertex(v *Vertex) (*Vertex, error) {
	g.ctx.Lock()
	defer g.ctx.Unlock()

	if g.ctx.visited {
		return nil, ErrDuplicateCompile
	}

	if v.cid > 0 {
		return g.vertex[v.cid-1], nil
	}

	// 初始化顶点基础上下文
	v.ctx.baseCtx, v.ctx.baseCancel = context.WithCancel(g.ctx.baseCtx)

	// 设置顶点id
	v.cid = int64(len(g.vertex) + 1)

	g.vertex = append(g.vertex, v)
	v.graph = g

	join(fmt.Sprintf(vertexPrefix, g.Name(), v.Name()), v)

	return v, nil
}

// 处理当前节点及其邻接节点的任务逻辑。首先处理当前节点的任务函数，并根据返回的错误情况进行相应的处理。
// 然后，遍历邻接节点，对于每个邻接节点，更新其依赖数量，如果依赖数量为 0，则启动一个新的 Goroutine 来处理该邻接节点。
// 这样可以实现图算法中节点之间的并发处理和依赖关系的维护。
func (g *Graph) runVertex(v *Vertex, errCh chan<- error) {
	defer func() {
		v.ctx.mainCancel()
		g.wg.Done()
		remove(fmt.Sprintf(vertexPrefix, g.Name(), v.Name()))
	}()

	// 图形级暂停
	if g.Paused() {
		select {
		case <-g.ctx.mainCtx.Done():
			// 被终止
			errCh <- ErrForceKill
			return
		case <-g.ctx.controlCtx.Done():
			// 继续
		}
	}

	// 节点级暂停
	if v.Paused() {
		select {
		case <-g.ctx.mainCtx.Done():
			// 被终止
			errCh <- ErrForceKill
			return
		case <-v.ctx.mainCtx.Done():
			// 被终止
			errCh <- ErrForceKill
			return
		case <-v.ctx.controlCtx.Done():
			// 继续
		}
	}

	// 设置运行中
	v.ctx.Lock()
	defer v.ctx.Unlock()
	v.running = true

	// 执行顶点函数
	if err := v.fn(v.ctx.mainCtx, g.Name(), v.Name()); err != nil {
		errCh <- err
		return
	}

	// 执行后面的顶点
	for k := range v.adjs {
		select {
		case <-g.ctx.mainCtx.Done():
			errCh <- ErrForceKill
			break
		default:
			dec := func() {
				g.ctx.Lock()
				defer g.ctx.Unlock()
				v.adjs[k].ndeps--
			}
			dec()
			if v.adjs[k].ndeps == 0 {
				g.wg.Add(1)
				go g.runVertex(v.adjs[k], errCh)
			}
		}
	}
}

// 合并多个上下文取消
func (g *Graph) withCancel(main, extra context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(main)
	stopf := context.AfterFunc(extra, func() {
		cancel()
	})
	return ctx, func() {
		cancel()
		stopf()
	}
}

// Run 运行任务流程,控制整个图算法的执行流程。它启动一个 Goroutine 来监听错误通道，同时遍历图的根节点并启动相应的 Goroutine 来处理每个根节点。
// 然后等待所有节点的处理完成，最后检查并打印可能发生的错误信息。同时，通过上下文的取消来通知所有 Goroutine 停止处理。
func (g *Graph) Run(ctx context.Context) error {
	if !g.ctx.visited {
		// 调用 compile 方法生成有向无环图，并检查是否有错误发生。
		// 如果有错误，则打印错误信息并返回。
		if err := g.compile(); err != nil {
			return err
		}
	}

	// 设置图形主上下文
	g.ctx.mainCtx, g.ctx.mainCancel = g.withCancel(ctx, g.ctx.baseCtx)

	// 设置所有顶点主上下文
	for k := range g.vertex {
		g.vertex[k].ctx.mainCtx, g.vertex[k].ctx.mainCancel = g.withCancel(g.ctx.mainCtx, g.vertex[k].ctx.baseCtx)
	}

	defer func() {
		g.ctx.mainCancel()
		remove(fmt.Sprintf(graphPrefix, g.Name()))
	}()

	var chError = make(chan error)
	var errs []error
	go func() {
		for err := range chError {
			errs = append(errs, err)
		}
	}()

	for k := range g.vertex {
		if !g.vertex[k].root {
			continue
		}
		g.wg.Add(1)
		go g.runVertex(g.vertex[k], chError)
	}

	g.wg.Wait()
	close(chError)

	if errs == nil {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

func (g *Graph) Validator() error {
	if g.vertex == nil {
		return ErrEmptyGraph
	}
	if g.ctx.visited {
		// 只允许编译一次图形
		// TODO: 多次编译
		return nil
	}
	return g.compile()
}

// compile 将任务列表转换为有向无环图。它遍历任务列表中的每个任务，创建对应的节点，并建立节点之间的依赖关系。
// 在建立依赖关系时，将依赖的节点的指针添加到当前节点的 adjs 切片中。
// 如果检测到回环，则返回 ErrCycleDetected 错误
func (g *Graph) compile() (err error) {
	g.ctx.Lock()
	defer func() {
		g.ctx.visited = true
		if err != nil {
			g.ctx.visited = false
			g.vertex = nil
		}
		g.ctx.Unlock()
	}()
	var nameMap = make(map[string]bool)

	// 循环遍历任务列表 g.vertex，为每个节点设置相应的属性，并根据任务的依赖关系将节点连接起来。
	for k := range g.vertex {
		if _, ok := nameMap[g.vertex[k].Name()]; ok {
			return ErrDuplicateVertexName
		}
		nameMap[g.vertex[k].Name()] = true

		g.vertex[k].ndeps = int64(len(g.vertex[k].deps))

		// 将当前顶点作为邻接或相邻分配给父顶点。
		// 具体地, 将依赖的节点的指针添加到当前节点的 adjs 切片中，表示当前节点依赖于这些节点。
		for _, dep := range g.vertex[k].deps {
			g.vertex[dep.cid-1].adjs = append(g.vertex[dep.cid-1].adjs, g.vertex[k])
		}

		// 如果任务没有依赖，将其节点添加到 roots 切片中。
		if len(g.vertex[k].deps) == 0 {
			g.vertex[k].root = true
		}

		// 检查图形是否存在回环, 检查以本节点为起点的子图是否存在回环。
		if err = g.detectCircularDependencies(g.vertex[k], []*Vertex{}); err != nil {
			return err
		}
	}
	return
}

// 使用深度优先搜索 (DFS) 的方式进行回环检测。它从给定的节点开始遍历邻接节点，并在遍历过程中检查是否存在回环。
// 如果发现已访问过的节点，则存在回环，返回 ErrCycleDetected 错误。否则，继续递归遍历邻接节点。
// 为了避免重复访问节点，使用 visited 属性对已访问的节点进行标记。
func (g *Graph) detectCircularDependencies(current *Vertex, path []*Vertex) error {
	// 如果发现某个邻接节点已经被访问过（即存在回环）
	if current.ctx.visited {
		return ErrCycleDetected
	}
	current.ctx.visited = true
	// 递归地遍历节点的邻接节点
	for k := range current.adjs {
		if err := g.detectCircularDependencies(current.adjs[k], append(path, current)); err != nil {
			return err
		}
	}

	current.ctx.visited = false
	return nil
}
