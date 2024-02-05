package dag

import (
	"context"
	"sync"
)

type iContext struct {
	sync.Mutex
	name    string // 名称
	visited bool   // 编译过或者追踪节点的遍历状态，以防止重复访问或陷入无限循环

	// 基础上下文. 控制未执行或执行中强杀
	baseCtx    context.Context
	baseCancel context.CancelFunc

	// 主上下文, 支持外部控制(超时,强杀)
	mainCtx    context.Context
	mainCancel context.CancelFunc

	// 控制上下文, 控制挂起或解卦
	controlCtx    context.Context
	controlCancel context.CancelFunc
}
