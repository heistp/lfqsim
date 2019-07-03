package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Tick uint64

type Seqno uint64

type Packet struct {
	Seqno     Seqno
	Timestamp Tick
	Size      int
	Hash      int
}

type FlowDef struct {
	Description      string
	Offset           Tick
	Interval         Tick
	IntervalVariance Tick
	Burst            int
	BurstVariance    int
	Size             int
	SizeVariance     int
}

type FlowState struct {
	NextEnqueue Tick
	NextSeqno   Seqno
	PriorSeqno  Seqno
}

type Config struct {
	EndTicks        Tick
	MTU             int
	MaxSize         int
	LateDump        bool
	LateDumpPackets bool
	Algorithm       string
	FlowDefs        []FlowDef
}

type FlowStats struct {
	BytesSent        uint64
	Throughput       float64
	MeanSojourn      float64
	MinSojourn       Tick
	MaxSojourn       Tick
	Enqueues         uint64
	Drops            uint64
	DropsPercent     float64
	SparseSends      uint64
	BulkSends        uint64
	TotalSends       uint64
	LateSends        uint64
	LateSendsPercent float64
	TotalSojourn     Tick `json:"-"`
}

type Results struct {
	FlowStats []FlowStats
}

type Queuer interface {
	Enqueue(p *Packet)

	Dequeue() (sent bool)

	Dump(reason string, packets bool)
}

type Sender interface {
	Send(p *Packet, sparse bool)
}

type Simulator struct {
	*Config
	FlowStates  []FlowState
	Results     Results
	Now         Tick
	NextDequeue Tick
	Queuer      Queuer
	Rand        *rand.Rand
}

func NewSimulator(c *Config) *Simulator {
	s := &Simulator{c,
		make([]FlowState, len(c.FlowDefs)),
		Results{
			make([]FlowStats, len(c.FlowDefs)),
		},
		0,
		0,
		nil,
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	if s.Algorithm == "CNQ" {
		s.Queuer = NewCNQ(len(c.FlowDefs), c.MaxSize, s)
	} else {
		s.Queuer = NewLFQ(len(c.FlowDefs), c.MaxSize, c.MTU, s)
	}
	return s
}

func (s *Simulator) Run() *Results {
	// initialize first enqueue time to flow offsets
	for i := 0; i < len(s.FlowDefs); i++ {
		s.FlowStates[i].NextEnqueue = s.FlowDefs[i].Offset
	}

	// run simulation
	for s.Now = 0; s.Now < s.EndTicks; s.Now++ {
		// call enqueue for each eligible flow
		for i := 0; i < len(s.FlowStates); i++ {
			fs := &s.FlowStates[i]
			r := &s.Results.FlowStats[i]
			if fs.NextEnqueue == s.Now {
				fd := &s.FlowDefs[i]
				for j := 0; j < fd.Burst+s.varyInt(fd.BurstVariance); j++ {
					s.Queuer.Enqueue(&Packet{fs.NextSeqno, s.Now, fd.Size + s.varyInt(fd.SizeVariance), i})
					r.Enqueues++
					fs.NextSeqno++
				}
				fs.NextEnqueue += fd.Interval + s.varyTick(fd.IntervalVariance)
			}
		}

		// call dequeue if it's time- and try again each tick if nothing was sent
		if s.Now == s.NextDequeue && !s.Queuer.Dequeue() {
			s.NextDequeue++
		}
	}

	// post-process results
	for i := 0; i < len(s.FlowDefs); i++ {
		r := &s.Results.FlowStats[i]
		r.TotalSends = r.SparseSends + r.BulkSends
		r.Throughput = 1000 * float64(r.BytesSent) / (float64(s.EndTicks) - float64(s.FlowDefs[i].Offset))
		r.MeanSojourn = float64(r.TotalSojourn) / float64(r.TotalSends)
		r.LateSendsPercent = float64(100.0) * float64(r.LateSends) / float64(r.TotalSends)
		r.Drops = r.Enqueues - r.TotalSends
		r.DropsPercent = float64(100.0) * float64(r.Drops) / float64(r.Enqueues)
	}

	return &s.Results
}

func (s *Simulator) Send(p *Packet, sparse bool) {
	i := p.Hash
	r := &s.Results.FlowStats[i]
	r.BytesSent += uint64(p.Size)
	sojourn := s.Now - p.Timestamp

	r.TotalSojourn += sojourn
	if sojourn < r.MinSojourn || r.MinSojourn == 0 {
		r.MinSojourn = sojourn
	}
	if sojourn > r.MaxSojourn {
		r.MaxSojourn = sojourn
	}

	if sparse {
		r.SparseSends++
	} else {
		r.BulkSends++
	}

	if p.Seqno < s.FlowStates[i].PriorSeqno {
		r.LateSends++
		if s.LateDump {
			s.Queuer.Dump(fmt.Sprintf("late packet %+v", p), s.LateDumpPackets)
		}
	}
	s.FlowStates[i].PriorSeqno = p.Seqno

	// schedule dequeue based on constant bitrate of one byte per tick
	s.NextDequeue += Tick(p.Size)
}

func (s *Simulator) varyInt(v int) int {
	if v == 0 {
		return 0
	}
	return s.Rand.Intn(2*v+1) - v
}

func (s *Simulator) varyTick(v Tick) Tick {
	if v == 0 {
		return 0
	}
	return Tick(s.Rand.Uint64())%(2*v+1) - v
}

// Additional simulation-specific methods for dumping state

func (q *LFQ) Dump(reason string, packets bool) {
	log.Printf("LFQ state dump (reason: %s):", reason)
	for i, bkt := range q.buckets {
		log.Printf("  Bucket %d: backlog=%d, deficit=%d, skip=%t", i, bkt.Backlog,
			bkt.Deficit, bkt.Skip)
	}
	q.Sparse.Dump("Sparse", packets)
	q.Bulk.Dump("Bulk", packets)
}

func (q *CNQ) Dump(reason string, packets bool) {
	log.Printf("CNQ state dump (reason: %s):", reason)
	for i, b := range q.backlogs {
		log.Printf("  Backlog %d: backlog=%d", i, b)
	}
	q.Sparse.Dump("Sparse", packets)
	q.Bulk.Dump("Bulk", packets)
}

func (q *Queue) Dump(label string, packets bool) {
	log.Printf("  Queue state (%s), Length: %d, Size: %d", label, q.Len(), q.Size)
	if packets {
		for i, p := range q.packets {
			log.Printf("%sPacket %d: %+v", "    ", i, p)
		}
	}
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
