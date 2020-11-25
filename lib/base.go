package lib

import "time"

type CallResult struct {

}

type RawReq struct {
	ID uint32
	Req []byte
}

type RawResp struct {
	ID     int64
	Resp   []byte
	Err    error
	Elapse time.Duration
}