package lib

import (
	"time"
)

type Caller interface {
	BuildReq() RawReq
	Call(RawReq, time.Duration)(RawResp, error)
	CheckRawResp(RawReq, RawResp) *CallResult
}
