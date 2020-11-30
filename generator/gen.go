package generator

import (
	"context"
	"fmt"
	"generator/lib"
	"math"
	"sync/atomic"
	"time"
)

type generator struct {
	timeoutNS   time.Duration
	lps         uint32
	durationNs  time.Duration
	caller      lib.Caller
	resultCh    chan *lib.CallResult
	ticketCount uint32
	ticketPool  lib.TicketPool
	cancelFunc  context.CancelFunc
}

func NewGenerator(param Param) (Generator,error){
	err := param.Check()
	if err != nil {
		return nil,err
	}
	gen := &generator{
		timeoutNS: param.TimeoutNS,
		lps: param.LPS,
		durationNs: param.DurationNS,
		caller: param.Caller,
		resultCh: param.ResultCh,
	}
	return gen,nil
}

func (gen *generator) init() {
	total64 := int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt64 {
		total64 = math.MaxInt64
	}
	gen.ticketCount = uint32(total64)
	gen.ticketPool = lib.NewTicketPool(gen.ticketCount)
}

func (gen *generator) Start() {
	if gen == nil {
		return
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), gen.durationNs)
	gen.cancelFunc = cancelFunc

	interval := time.Duration(1e9 / gen.lps)
	throttle := time.Tick(interval)
	gen.genLoad(ctx, throttle)
}

func (gen *generator) Stop() {
	gen.cancelFunc()
}

func (gen *generator) genLoad(ctx context.Context, throttle <-chan time.Time) {
	for {
		select {
		case <- ctx.Done():
			return
		default:
		}

		gen.asyncCall(ctx)

		if throttle != nil {
			select {
			case <- throttle:
			case <- ctx.Done():
				return
			}
		}
	}
}

func (gen *generator) asyncCall(ctx context.Context){
	gen.ticketPool.Take()
	go func(){
		defer gen.ticketPool.Return()

		rawReq := gen.caller.BuildReq()
		var callStatus uint32 = 0

		//超时处理
		timer := time.AfterFunc(gen.durationNs,
			func(){
				if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
					return
				}	
				result := &lib.CallResult{
					ID:     rawReq.ID,
					Req:    rawReq,
					Code:   lib.RET_CODE_WARNING_CALL_TIMEOUT,
					Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
					Elapse: gen.timeoutNS,
				}
				gen.sendResult(result)			
			},
		)
		startTime := time.Now().UnixNano()
		rawResp, err := gen.caller.Call(rawReq,gen.durationNs)
		endTime := time.Now().UnixNano()
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			return
		}			
		timer.Stop()
		elapse:=time.Duration(endTime-startTime)

		if err != nil {
			result := &lib.CallResult{
				ID:     rawReq.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_ERROR_CALL,
				Msg:    fmt.Sprintf("call fail: %s)", err),
				Elapse: elapse,
			}
			gen.sendResult(result)	
		}

		result := gen.caller.CheckRawResp(rawReq, rawResp)	
		gen.sendResult(result)
	}()			
}

func (gen *generator) sendResult(result *lib.CallResult){
	select {
	case gen.resultCh <- result:
	default: //通道已满日志
	}
}



