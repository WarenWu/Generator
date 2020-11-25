package generator

import(
	"generator/lib"
	"time"
)
type Param struct {
	caller lib.Caller
	resultCh chan *lib.Resultb
	timeoutNS  time.Duration
	lps uint32
	durationNs time.Duration
}

func (*Param)Check() bool {
	return true
}