package main

import "log"

// Differences from I-D pseudo-code:
// - No AQM
// - Packet hash specified directly, so no cached value needed
// - Send method contains sparse flag for simulation stats
// - Pull method contains fast flag for experimental fast pull
// - Timestamp is a Tick for the simulation

// Algorithm / I-D notes:
// - Walking all buckets in dequeue may be expensive
// - Enqueue might loop infinitely if MaxSize too small. What's the minimum?

// Todo:
// - record dropped packets
// - dump state on late packet

type Packet struct {
	Seqno     uint64
	Timestamp Tick
	Size      int
	Hash      int
}

type FlowBucket struct {
	Backlog int
	Deficit int
	Skip    bool
}

type Queue struct {
	packets   []*Packet
	ScanIndex int
}

func NewQueue() *Queue {
	return &Queue{make([]*Packet, 0), 0}
}

func (q *Queue) Len() int {
	return len(q.packets)
}

func (q *Queue) Size() (s int) {
	for i := 0; i < len(q.packets); i++ {
		s += q.packets[i].Size
	}
	return
}

func (q *Queue) Pop() (p *Packet, found bool) {
	if len(q.packets) > 0 {
		p, q.packets = q.packets[0], q.packets[1:]
		found = true
	}
	return
}

func (q *Queue) Push(p *Packet) {
	q.packets = append(q.packets, p)
	return
}

func (q *Queue) Head() (p *Packet, found bool) {
	if len(q.packets) > 0 {
		p = q.packets[0]
		found = true
	}
	return
}

func (q *Queue) Empty() bool {
	return len(q.packets) == 0
}

func (q *Queue) Scan() (p *Packet, found bool) {
	if q.ScanIndex < len(q.packets) {
		p = q.packets[q.ScanIndex]
		found = true
	}
	return
}

func (q *Queue) Pull(fast bool) (p *Packet, found bool) {
	if q.ScanIndex < len(q.packets) {
		p = q.packets[q.ScanIndex]
		found = true
		if fast {
			q.packets[q.ScanIndex] = q.packets[len(q.packets)-1]
			q.packets = q.packets[:len(q.packets)-1]
		} else {
			q.packets = append(q.packets[:q.ScanIndex], q.packets[q.ScanIndex+1:]...)
		}
	}
	return
}

type Sender interface {
	Send(p *Packet, sparse bool)
}

type LFQ struct {
	Sparse   *Queue
	Bulk     *Queue
	buckets  []FlowBucket
	MaxSize  int
	MTU      int
	FastPull bool
	Sender   Sender
}

func NewLFQ(maxFlows int, maxSize int, MTU int, fastPull bool, s Sender) *LFQ {
	return &LFQ{
		NewQueue(),
		NewQueue(),
		make([]FlowBucket, maxFlows),
		maxSize,
		MTU,
		fastPull,
		s,
	}
}

func (q *LFQ) Enqueue(p *Packet, t Tick) {
	for q.Sparse.Size()+q.Bulk.Size()+p.Size > q.MaxSize {
		// queue overflow, drop from bulk head
		if dp, found := q.Bulk.Pop(); found {
			q.buckets[dp.Hash].Backlog -= 1
		} else {
			// avoid infinite loop if MaxSize too small
			log.Println("lfq: avoided infinite loop in enqueue")
			break
		}
	}

	bkt := &q.buckets[p.Hash]
	p.Timestamp = t

	if bkt.Backlog == 0 && bkt.Deficit >= 0 && !bkt.Skip {
		q.Sparse.Push(p)
	} else {
		q.Bulk.Push(p)
	}
	bkt.Backlog++
}

func (q *LFQ) Dequeue() {
	var p *Packet
	var found bool

	// Sparse queue gets strict priority
	p, found = q.Sparse.Pop()
	if found {
		q.Sender.Send(p, true)
		bkt := &q.buckets[p.Hash]
		q.sent(p, bkt)
		return
	}

	// Process Bulk queue if Sparse queue was empty
	for !q.Bulk.Empty() {
		if p, found = q.Bulk.Scan(); !found {
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
			p, _ = q.Bulk.Scan()
		}

		if bkt := &q.buckets[p.Hash]; !bkt.Skip {
			// packet eligible for immediate delivery
			q.Sender.Send(p, false)
			q.Bulk.Pull(q.FastPull)
			q.sent(p, bkt)
			return
		} else {
			// packet stays in queue
			q.Bulk.ScanIndex++
		}
	}
}

func (q *LFQ) sent(p *Packet, bkt *FlowBucket) {
	bkt.Backlog--
	bkt.Deficit -= p.Size
	if bkt.Deficit < 0 {
		bkt.Skip = true
		bkt.Deficit += q.MTU
	}
}
