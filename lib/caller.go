package lib

import (
	"gopcp.v2/chapter4/loadgen/lib"
	"time"
)

type Caller interface {
	BuildReq() RawReq
	Call(RawQeq, timeoutNS time.Duration)([]byte,error)
	CheckRawResp(RawResp) *lib.CallResult
}
