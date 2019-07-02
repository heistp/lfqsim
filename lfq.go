package main

type FlowBucket struct {
	Backlog int
	Deficit int
	Skip    bool
}

type LFQ struct {
	Sparse  *Queue
	Bulk    *ScanQueue
	buckets []FlowBucket
	MaxSize int
	MTU     int
	Sender  Sender
}

func NewLFQ(maxFlows int, maxSize int, MTU int, s Sender) *LFQ {
	return &LFQ{
		NewQueue(),
		NewScanQueue(),
		make([]FlowBucket, maxFlows),
		maxSize,
		MTU,
		s,
	}
}

func (q *LFQ) Enqueue(p *Packet) {
	for q.Sparse.Size+q.Bulk.Size+p.Size > q.MaxSize {
		// queue overflow, drop first from bulk head, then from sparse
		dp := q.Bulk.Pop()
		if dp == nil {
			dp = q.Sparse.Pop()
		}
		q.buckets[dp.Hash].Backlog -= 1
	}

	bkt := &q.buckets[p.Hash]

	if bkt.Backlog == 0 && bkt.Deficit >= 0 && !bkt.Skip {
		q.Sparse.Push(p)
	} else {
		q.Bulk.Push(p)
	}
	bkt.Backlog++
}

func (q *LFQ) Dequeue() (sent bool) {
	var p *Packet

	// Sparse queue gets strict priority
	if p = q.Sparse.Pop(); p != nil {
		q.Sender.Send(p, true)
		bkt := &q.buckets[p.Hash]
		q.sent(p, bkt)
		sent = true
		return
	}

	// Process Bulk queue only if Sparse queue was empty
	for !q.Bulk.Empty() {
		if p = q.Bulk.Scan(); p == nil {
			// scan has reached tail of queue
			for i := 0; i < len(q.buckets); i++ {
				bkt := &q.buckets[i]
				if !bkt.Skip {
					if bkt.Backlog == 0 {
						bkt.Deficit = 0
					}
				} else {
					bkt.Skip = false
				}
			}

			q.Bulk.ScanIndex = 0
			p = q.Bulk.Scan()
		}

		if bkt := &q.buckets[p.Hash]; !bkt.Skip {
			// packet eligible for immediate delivery
			q.Sender.Send(p, false)
			q.Bulk.Pull()
			q.sent(p, bkt)
			sent = true
			return
		} else {
			// packet stays in queue
			q.Bulk.ScanIndex++
		}
	}

	return
}

func (q *LFQ) sent(p *Packet, bkt *FlowBucket) {
	bkt.Backlog--
	bkt.Deficit -= p.Size
	if bkt.Deficit < 0 {
		bkt.Skip = true
		bkt.Deficit += q.MTU
	}
}
