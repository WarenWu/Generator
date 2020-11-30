package generator

import (
	"context"
	"generator/lib"
	"math"

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


func GetGenerator(param Param) Generator {
	if !param.Check() {
		//打印错误日志
	}
	gen := &generator{
		timeoutNS:  param.timeoutNS,
		lps:        param.lps,
		caller:     param.caller,
		resultCh:   param.resultCh,
	}
	gen.init()
	return gen
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
	if(gen == nil){
		//错误日志
	}
	ctx, cancelFunc = context.WithTimeout(context.Background(), gen.durationNS)
	gen.cancelFunc = cancelFunc

	//var throttle = <-chan time.Time
	if(gen.lps > 0){
		interval := time.Duration(1e9 / gen.lps)
		throttle = time.Tick(interval)
		genLoad(ctx, throttle)
	}else if (gen.lps) {
		interval := 0
		throttle = time.Tick(interval)
		genLoad(ctx, throttle)
	}else {
		//错误日志
	}

	
}

func (gen *generator) Stop() {
	gen.cancelFunc()
}

func (gen *generator) genLoad(ctx context.Context, throttle <-chan time.Time) {
	for {
		select {
		case <- gen.ctx.Done():
			return
		default:
		}

		go gen.asyncCall(ctx)

		if(throttle != nil){
			select {
			case <- throttle:
			case <- gen.ctx.Done():
				return
			}
		}
				
	}
}

func (gen *generator) asyncCall(ctx context.Context){
		
	go func(){
		gen.ticketPool.Take()

		rawReq := lib.BuildReq()
		var callStatus uint32 = 0;

		//超时处理
		timer ：= time.AfterFunc(gen.durationNs,
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
			}
		)
		starTime := time.Now.UnixNano
		resp, err := gen.caller.Call(ctx, rawReq)
		endTime := time.Now.UnixNano
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			return
		}			
		timer.Stop()

		if(err != nil){
			result := &lib.CallResult{
				ID:     rawReq.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_ERROR_CALL,
				Err:	err
				Msg:    fmt.Sprintf("call fail: %s)", err),
				Elapse: time.Duration(endTime - starTime),
			}
			gen.sendResult(result)	
		}
		
		rawResp := lib.RawResp{
			ID:rawReq.ID,
			Resp:resp,
			Err:err,
			Elapse:time.Duration(endTime - starTime),
		}
		result := gen.caller.CheckRawResp(rawReq, rawResp)	
		gen.sendResult(result)
	}()			
}

func (gen *generator) sendResult(result *lib.CallResult){
	select {
	case gen.resultCh <- result:
	case default:
		//通道已满日志
	}
}



