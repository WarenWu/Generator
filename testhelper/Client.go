package testhelper

import (
	"encoding/json"
	"fmt"
	"generator/lib"
	"math"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

type TcpClient struct {
	addr string
}

func NewTcpClient(addr string) lib.Caller {
	var TcpClient = TcpClient{addr}
	return &TcpClient
}

func (client *TcpClient) BuildReq() (rawReq lib.RawReq) {
	id := time.Now().UnixNano()
	req := Request{
		ID: id,
		Operands: []int32{
			int32(rand.Int31n(math.MaxInt32) + 1),
			int32(rand.Int31n(math.MaxInt32) + 1)},
		Operator: func() string {
			return operators[rand.Int31n(100)%4]
		}(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {

	}
	rawReq.Req = reqBytes
	return
}

func (client *TcpClient) Call(rawReq lib.RawReq, timeout time.Duration) (lib.RawResp, error) {
	var flag uint32
	var rawResp lib.RawResp
	var err error
	rawResp.ID = rawReq.ID

	conn, err := net.DialTimeout("tcp", client.addr, timeout)
	if err != nil {
		rawResp.Err = err
	}

	err = write(conn, rawReq.Req)
	if err != nil {
		rawResp.Err = err
	}

	readBytes, err := read(conn)
	if err != nil {
		rawResp.Err = err
	}

	if atomic.CompareAndSwapUint32(&flag, 0, 1) {
		rawResp.Err = err
		rawResp.Resp = readBytes
	}

	return rawResp, err
}

func (client *TcpClient) CheckRawResp(rawReq lib.RawReq, rawResp lib.RawResp) *lib.CallResult {
	var result lib.CallResult
	result.ID = rawReq.ID
	result.Req = rawReq
	result.Resp = rawResp

	var req Request
	err := json.Unmarshal(rawReq.Req, &req)
	if err != nil {
		result.Code = lib.RET_CODE_FATAL_CALL
		result.Msg =
			fmt.Sprintf("Incorrectly formatted Req: %s!\n", string(rawReq.Req))
		return &result
	}
	var resp Response
	err = json.Unmarshal(rawResp.Resp, &resp)
	if err != nil {
		result.Code = lib.RET_CODE_ERROR_RESPONSE
		result.Msg =
			fmt.Sprintf("Incorrectly formatted Resp: %s!\n", string(rawResp.Resp))
		return &result
	}
	if resp.ID != req.ID {
		result.Code = lib.RET_CODE_ERROR_RESPONSE
		result.Msg =
			fmt.Sprintf("Inconsistent raw id! (%d != %d)\n", rawReq.ID, rawResp.ID)
		return &result
	}
	if resp.Err != nil {
		result.Code = lib.RET_CODE_ERROR_CALEE
		result.Msg =
			fmt.Sprintf("Abnormal server: %s!\n", resp.Err)
		return &result
	}
	if resp.Result != op(req.Operands, req.Operator) {
		result.Code = lib.RET_CODE_ERROR_RESPONSE
		result.Msg =
			fmt.Sprintf(
				"Incorrect result: %s!\n",
				genFormula(req.Operands, req.Operator, resp.Result, false))
		return &result
	}
	result.Code = lib.RET_CODE_SUCCESS
	result.Msg = fmt.Sprintf("Success. (%s)", resp.Formula)
	return &result
}
