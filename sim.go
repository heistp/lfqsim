package main

import (
	"math/rand"
	"time"
)

type Tick uint64

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
	NextSeqno   uint64
	PriorSeqno  uint64
}

type Config struct {
	EndTicks        Tick
	DequeueInterval Tick
	MTU             int
	FastPull        bool
	MaxSize         int
	FlowDefs        []FlowDef
}

type FlowStats struct {
	BytesSent        uint64
	Throughput       float64
	MeanSojourn      float64
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

type Simulator struct {
	*Config
	FlowStates []FlowState
	Results    Results
	Tick       Tick
	Rand       *rand.Rand
}

func NewSimulator(c *Config) *Simulator {
	s := &Simulator{c,
		make([]FlowState, len(c.FlowDefs)),
		Results{
			make([]FlowStats, len(c.FlowDefs)),
		},
		0,
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for i := 0; i < len(c.FlowDefs); i++ {
		s.FlowStates[i].NextEnqueue = c.FlowDefs[i].Offset
	}
	return s
}

func (s *Simulator) Send(p *Packet, sparse bool) {
	i := p.Hash
	r := &s.Results.FlowStats[i]
	r.BytesSent += uint64(p.Size)
	r.TotalSojourn += (s.Tick - p.Timestamp)

	if sparse {
		r.SparseSends++
	} else {
		r.BulkSends++
	}

	if p.Seqno < s.FlowStates[i].PriorSeqno {
		r.LateSends++
	}
	s.FlowStates[i].PriorSeqno = p.Seqno
}

func (s *Simulator) Run() *Results {
	q := NewLFQ(len(s.FlowDefs), s.MaxSize, s.MTU, s.FastPull, s)

	// run simulation
	for s.Tick = 0; s.Tick < s.EndTicks; s.Tick++ {
		// call enqueue for each eligible flow
		for i := 0; i < len(s.FlowStates); i++ {
			fs := &s.FlowStates[i]
			r := &s.Results.FlowStats[i]
			if fs.NextEnqueue == s.Tick {
				fd := &s.FlowDefs[i]
				for j := 0; j < fd.Burst+s.randVaryInt(fd.BurstVariance); j++ {
					q.Enqueue(&Packet{fs.NextSeqno, 0, fd.Size + s.randVaryInt(fd.SizeVariance), i}, s.Tick)
					r.Enqueues++
					fs.NextSeqno++
				}
				fs.NextEnqueue += fd.Interval + s.randVaryTick(fd.IntervalVariance)
			}
		}

		// call dequeue
		if s.Tick%s.DequeueInterval == 0 {
			q.Dequeue()
		}
	}

	// post-process results
	for i := 0; i < len(s.FlowDefs); i++ {
		r := &s.Results.FlowStats[i]
		r.TotalSends = r.SparseSends + r.BulkSends
		r.Throughput = float64(s.DequeueInterval) * float64(r.BytesSent) /
			(float64(s.EndTicks) - float64(s.FlowDefs[i].Offset))
		r.MeanSojourn = float64(r.TotalSojourn) / float64(r.TotalSends)
		r.LateSendsPercent = float64(100.0) * float64(r.LateSends) / float64(r.TotalSends)
		r.Drops = r.Enqueues - r.TotalSends
		r.DropsPercent = float64(100.0) * float64(r.Drops) / float64(r.Enqueues)
	}

	return &s.Results
}

func (s *Simulator) randVaryInt(v int) int {
	if v == 0 {
		return 0
	}
	return s.Rand.Intn(2*v+1) - v
}

func (s *Simulator) randVaryTick(v Tick) Tick {
	if v == 0 {
		return 0
	}
	return Tick(s.Rand.Uint64())%(2*v+1) - v
}
