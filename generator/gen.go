package generator

import (
	"generator/lib"
	"math"
	"time"
)

type generator struct {
	timeoutNS  time.Duration
	lps uint32
	durationNs time.Duration
	caller lib.Caller
	resultCh chan *lib.CallResult
	ticketCount uint32
	ticketPool  lib.TicketPool
}

func (this *generator)Start(){

}

func (this *generator)Stop(){

}

func NewGenerator(param Param) (Generator,error){
	if!param.Check(){
		//打印错误日志
		return nil,nil
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

func (this *generator) init() {
	total64 := int64(this.timeoutNS)/int64(1e9/this.lps) + 1
	if total64 > math.MaxInt64{
		total64 = math.MaxInt64
	}
	this.ticketCount = uint32(total64)
	this.ticketPool = lib.NewTicketPool(this.ticketCount)
}

func (this *generator) callOne() {

}