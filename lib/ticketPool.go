package lib

type TicketPool interface {

	// 拿走一张票。
	Take()
	// 归还一张票。
	Return()
	// 票池是否已被激活。
	Active() bool
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
	myPool := &ticketPool{
		total: total,
		active: true,
	}
	ch := make(chan struct{}, total)
	myPool.poolCh = ch

	for i:=0; i < int(total); i++{
		ch <- struct{}{}
	}
	return myPool
}


func (this *ticketPool)Take()  {

}

func (this *ticketPool)Return()  {

}

func (this *ticketPool)Active() bool {

}

func (this *ticketPool)Total() uint32 {

}

func (this *ticketPool)Remainder() uint32 {

}




