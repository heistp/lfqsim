package main

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

type ScanQueue struct {
	*Queue
	ScanIndex int
}

func NewScanQueue() *ScanQueue {
	return &ScanQueue{NewQueue(), 0}
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
		q.packets = append(q.packets[:q.ScanIndex], q.packets[q.ScanIndex+1:]...)
		q.Size -= p.Size
	}
	return
}

func (q *ScanQueue) Empty() bool {
	return len(q.packets) == 0
}
