package main

// Cheap Nasty Queueing
type CNQ struct {
	Sparse   *Queue
	Bulk     *Queue
	backlogs []int
	MaxSize  int
	Sender   Sender
}

func NewCNQ(maxFlows int, maxSize int, s Sender) *CNQ {
	return &CNQ{
		NewQueue(),
		NewQueue(),
		make([]int, maxFlows),
		maxSize,
		s,
	}
}

func (q *CNQ) Enqueue(p *Packet) {
	// queue overflow, drop first from bulk
	for q.Sparse.Size+q.Bulk.Size+p.Size > q.MaxSize {
		dp := q.Bulk.Pop()
		if dp == nil {
			break
		}
		q.backlogs[dp.Hash]--
	}

	// then drop from sparse if still needed
	for q.Sparse.Size+q.Bulk.Size+p.Size > q.MaxSize {
		dp := q.Sparse.Pop()
		if dp == nil {
			break
		}
		q.backlogs[dp.Hash]--
	}

	if q.backlogs[p.Hash] == 0 {
		q.Sparse.Push(p)
		q.Bulk.Push(&Packet{p.Seqno, p.Timestamp, 0, p.Hash})
		q.backlogs[p.Hash] = 2
	} else {
		q.Bulk.Push(p)
		q.backlogs[p.Hash]++
	}
}

func (q *CNQ) Dequeue() (sent bool) {
	var p *Packet

	// Sparse queue gets strict priority
	if p = q.Sparse.Pop(); p != nil {
		q.Sender.Send(p, true)
		q.backlogs[p.Hash]--
		sent = true
		return
	}

	// Process Bulk queue only if Sparse queue was empty
	for p = q.Bulk.Pop(); p != nil; p = q.Bulk.Pop() {
		if p.Size > 0 {
			q.Sender.Send(p, false)
			q.backlogs[p.Hash]--
			sent = true
			return
		} else {
			q.backlogs[p.Hash]--
		}
	}

	return
}
