package lib

import (
	"context"

)

type Caller interface {
	BuildReq() RawReq
	Call(ctx context.Context, rawReq RawReq) ([]byte, error)
	CheckRawResp(rawReq RawReq, rawResp RawResp) *CallResult
}
