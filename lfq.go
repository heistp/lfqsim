package main

import "log"

// Differences from I-D pseudo-code:
// - No AQM
// - Timestamp is a Tick for the simulation
// - Packet hash specified directly, so no cached value needed
// - FastPull flag uses experimental fast pull (causes packet re-ordering)
// - Send method contains sparse flag for simulation stats

// Algorithm / I-D notes:
// - Enqueue might loop infinitely if MaxSize is less than the size of the
//   sparse queue + the current packet being enqueued

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
	packets []*Packet
	Size    int
}

func NewQueue() *Queue {
	return &Queue{make([]*Packet, 0), 0}
}

func (q *Queue) Len() int {
	return len(q.packets)
}

func (q *Queue) Pop() (p *Packet) {
	if len(q.packets) > 0 {
		p, q.packets = q.packets[0], q.packets[1:]
		q.Size -= p.Size
	}
	return
}

func (q *Queue) Push(p *Packet) {
	q.packets = append(q.packets, p)
	q.Size += p.Size
	return
}

func (q *Queue) Dump(label string, packets bool) {
	log.Printf("  Queue state (%s), Length: %d, Size: %d", label, q.Len(), q.Size)
	if packets {
		for i, p := range q.packets {
			log.Printf("%sPacket %d: %+v", "    ", i, p)
		}
	}
}

type ScanQueue struct {
	*Queue
	ScanIndex int
	FastPull  bool
}

func NewScanQueue(fastPull bool) *ScanQueue {
	return &ScanQueue{NewQueue(), 0, fastPull}
}

func (q *ScanQueue) Scan() (p *Packet) {
	if q.ScanIndex < len(q.packets) {
		p = q.packets[q.ScanIndex]
	}
	return
}

func (q *ScanQueue) Pop() (p *Packet) {
	p = q.Queue.Pop()
	if p != nil && q.ScanIndex > 0 {
		q.ScanIndex--
	}
	return
}

func (q *ScanQueue) Pull() (p *Packet) {
	if q.ScanIndex < len(q.packets) {
		p = q.packets[q.ScanIndex]
		if q.FastPull {
			q.packets[q.ScanIndex] = q.packets[len(q.packets)-1]
			q.packets = q.packets[:len(q.packets)-1]
		} else {
			q.packets = append(q.packets[:q.ScanIndex], q.packets[q.ScanIndex+1:]...)
		}
		q.Size -= p.Size
	}
	return
}

func (q *ScanQueue) Empty() bool {
	return len(q.packets) == 0
}

func (q *ScanQueue) Dump(label string, packets bool) {
	log.Printf("  Queue state (%s), Length: %d, Size: %d, ScanIndex: %d",
		label, q.Len(), q.Size, q.ScanIndex)
	if packets {
		for i, p := range q.packets {
			var prefix string
			if i == q.ScanIndex {
				prefix = "  ->"
			} else {
				prefix = "    "
			}
			log.Printf("%sPacket %d: %+v", prefix, i, p)
		}
	}
}

type Sender interface {
	Send(p *Packet, sparse bool, q *LFQ)
}

type LFQ struct {
	Sparse  *Queue
	Bulk    *ScanQueue
	buckets []FlowBucket
	MaxSize int
	MTU     int
	Sender  Sender
}

func NewLFQ(maxFlows int, maxSize int, MTU int, fastPull bool, s Sender) *LFQ {
	return &LFQ{
		NewQueue(),
		NewScanQueue(fastPull),
		make([]FlowBucket, maxFlows),
		maxSize,
		MTU,
		s,
	}
}

func (q *LFQ) Enqueue(p *Packet, t Tick) {
	for q.Sparse.Size+q.Bulk.Size+p.Size > q.MaxSize {
		// queue overflow, drop from bulk head
		if dp := q.Bulk.Pop(); dp != nil {
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

func (q *LFQ) Dequeue() (sent bool) {
	var p *Packet

	// Sparse queue gets strict priority
	if p = q.Sparse.Pop(); p != nil {
		q.Sender.Send(p, true, q)
		bkt := &q.buckets[p.Hash]
		q.sent(p, bkt)
		sent = true
		return
	}

	// Process Bulk queue if Sparse queue was empty
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
			q.Sender.Send(p, false, q)
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

func (q *LFQ) Dump(label string, packets bool) {
	log.Printf("LFQ state dump (reason: %s):", label)
	for i, bkt := range q.buckets {
		log.Printf("  Bucket %d: backlog=%d, deficit=%d, skip=%t", i, bkt.Backlog,
			bkt.Deficit, bkt.Skip)
	}
	q.Sparse.Dump("Sparse", packets)
	q.Bulk.Dump("Bulk", packets)
}
