package lib

type TicketPool interface {

	// 拿走一张票。
	Take()
	// 归还一张票。
	Return()
	// 票的总数。
	Total() uint32
	// 剩余的票数。
	Remainder() uint32
}

type ticketPool struct {
	total uint32
	poolCh chan struct{}
	active bool
}

func NewTicketPool(total uint32)  TicketPool{
	myPool := ticketPool{
		total: total,
	}
	ch := make(chan struct{}, total)
	myPool.poolCh = ch

	for i:=0; i < int(total); i++{
		ch <- struct{}{}
	}
	return &myPool
}


func (t *ticketPool)Take()  {
	<-t.poolCh
}

func (t *ticketPool)Return()  {
	t.poolCh <- struct{}{}
}

func (t *ticketPool)Total() uint32 {
	return t.total
}

func (t *ticketPool)Remainder() uint32 {
	return uint32(len(t.poolCh))
}




