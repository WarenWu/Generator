package generator

import(
	"generator/lib"
	"time"
)
type Param struct {
	Caller     lib.Caller           // 调用器。
	TimeoutNS  time.Duration        // 响应超时时间，单位：纳秒。
	LPS        uint32               // 每秒载荷量。
	DurationNS time.Duration        // 负载持续时间，单位：纳秒。
	ResultCh   chan *lib.CallResult // 调用结果通道。
}

func (*Param)Check() bool {
	return true
}