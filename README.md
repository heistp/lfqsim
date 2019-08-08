# lfqsim

This is a discrete time simulator for Lightweight Fair Queueing (LFQ) and
Cheap Nasty Queueing (CNQ).

## Quick Start

`lfqsim` reads its configuration from stdin and writes results to stdout, both in
JSON format. After having installed Go using your package manager or from
[here](https://golang.org/dl/):

```
go get github.com/heistp/lfqsim
cd ~/go/src/github.com/heistp/lfqsim # or location of lfqsim directory
go build
./lfqsim < config-lfq.json
```

## Description

Lightweight Fair Queueing (LFQ) and Cheap Nasty Queueing (CNQ) are fair queueing
algorithms from Jonathon Morton with a small codebase, low memory requirements,
and a design suitable for implementation in hardware. While LFQ provides
throughput fairness, CNQ is even simpler and does not, providing only sparse
flow optimization.

This related project provides an LFQ implementation (in `lfq.go`), a CNQ
implementation (in `cnq.go`) and a discrete time simulation that measures how
they perform with a configurable number of flows and their characteristics.

Each unit of time is referred to as a tick. Enqueue is called for each
configured flow at regular or varying intervals in ticks. Dequeue is called at a
constant bitrate of one tick per byte. See the [Results](#results) section for
information on what results are obtained. Note that late packets (out of order
packets) are not expected, so please file an issue if you discover simulation
parameters that cause them, or any other pathological behavior.

## Differences from LFQ Internet Draft Specification

- No AQM (Active Queue Management), so we rely on overflow to manage queue size
- The packet hash is specified directly
- For simulation purposes:
  - Timestamp is a Tick
  - Send method contains sparse flag for stats

## Sample Simulation and Discussion

Using the default `config-lfq.json`, output similar to what's below may be produced.

```
{
    "EndTicks": 100000000,
    "MTU": 1500,
    "MaxSize": 256000,
    "LateDump": true,
    "LateDumpPackets": false,
    "Algorithm": "LFQ",
    "FlowDefs": [
        {
            "Description": "1200 byte packets at almost 1/3 interface rate",
            "Offset": 0,
            "Interval": 3700,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1200,
            "SizeVariance": 99
        },
        {
            "Description": "Aggressive flow- 1500 byte packets at 3x interface rate",
            "Offset": 0,
            "Interval": 500,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
            "Description": "1500 byte packets exactly at interface rate",
            "Offset": 0,
            "Interval": 1500,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
            "Description": "Bursty flow- 1500 byte packets with 32 packet bursts at interface rate",
            "Offset": 0,
            "Interval": 48000,
            "IntervalVariance": 0,
            "Burst": 32,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
            "Description": "Smallish packets- 200 byte packets, burst 2 +/- 2, at interval 1000",
            "Offset": 0,
            "Interval": 10000,
            "IntervalVariance": 50,
            "Burst": 2,
            "BurstVariance": 2,
            "Size": 200,
            "SizeVariance": 50
        },
        {
            "Description": "Small packets- 50 byte packets at interval 750",
            "Offset": 0,
            "Interval": 750,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 50,
            "SizeVariance": 10
        },
        {
            "Description": "Sparse flow- 100 byte packets at interval 20000",
            "Offset": 0,
            "Interval": 20000,
            "IntervalVariance": 100,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 100,
            "SizeVariance": 0
        }
    ]
}
{
    "FlowStats": [
        {
            "BytesSent": 22464062,
            "Throughput": 224.64062,
            "MeanSojourn": 46938.40334473178,
            "MinSojourn": 7624,
            "MaxSojourn": 50892,
            "Enqueues": 27028,
            "Drops": 8312,
            "DropsPercent": 30.753292881456268,
            "SparseSends": 1,
            "BulkSends": 18715,
            "TotalSends": 18716,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 22465500,
            "Throughput": 224.655,
            "MeanSojourn": 48908.11771382787,
            "MinSojourn": 1180,
            "MaxSojourn": 50895,
            "Enqueues": 200000,
            "Drops": 185023,
            "DropsPercent": 92.5115,
            "SparseSends": 1,
            "BulkSends": 14976,
            "TotalSends": 14977,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 22465500,
            "Throughput": 224.655,
            "MeanSojourn": 48671.72658075716,
            "MinSojourn": 2680,
            "MaxSojourn": 50995,
            "Enqueues": 66667,
            "Drops": 51690,
            "DropsPercent": 77.53461232693836,
            "SparseSends": 1,
            "BulkSends": 14976,
            "TotalSends": 14977,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 22464000,
            "Throughput": 224.64,
            "MeanSojourn": 24438.113848824785,
            "MinSojourn": 30,
            "MaxSojourn": 51906,
            "Enqueues": 66688,
            "Drops": 51712,
            "DropsPercent": 77.54318618042227,
            "SparseSends": 1,
            "BulkSends": 14975,
            "TotalSends": 14976,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 2976785,
            "Throughput": 29.76785,
            "MeanSojourn": 3866.0304108485498,
            "MinSojourn": 1,
            "MaxSojourn": 14186,
            "Enqueues": 14896,
            "Drops": 0,
            "DropsPercent": 0,
            "SparseSends": 6449,
            "BulkSends": 8447,
            "TotalSends": 14896,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 6664774,
            "Throughput": 66.64774,
            "MeanSojourn": 5477.628905722643,
            "MinSojourn": 30,
            "MaxSojourn": 13288,
            "Enqueues": 133334,
            "Drops": 4,
            "DropsPercent": 0.0029999850000749996,
            "SparseSends": 8048,
            "BulkSends": 125282,
            "TotalSends": 133330,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 500000,
            "Throughput": 5,
            "MeanSojourn": 656.9394,
            "MinSojourn": 1,
            "MaxSojourn": 5962,
            "Enqueues": 5000,
            "Drops": 0,
            "DropsPercent": 0,
            "SparseSends": 5000,
            "BulkSends": 0,
            "TotalSends": 5000,
            "LateSends": 0,
            "LateSendsPercent": 0
        }
    ]
}
```

The output contains two objects, the parsed configuration, for verification, and
the results. From these LFQ results, it can be seen that throughput fairness is
achieved for the first four flows, which send at a sufficiently high rate. It
can also be seen that the sojourn time for the sparse flow (last flow) is far
less than that for the bulk flows, as expected. Flows with intermediate
characteristics can also be seen.

In a real implementation, AQM would manage the bulk queue, and flows would
respond to congestion signals, reducing their send rate. In this simulation,
that is ignored, and bulk flows may experience high drop rates in the
enforcement of throughput fairness.

Note that if the `MaxSize` parameter is too low, throughput fairness may
not be achieved, and bursty flows may see throughput cuts.

## Configuration

The JSON configuration object passed to stdin controls the simulation, and has
the following format:

* `EndTicks` the last tick for the simulation, exclusive
* `MTU` the MTU used for LFQ
* `MaxSize` LFQ's maximum size for the queue, in bytes
* `LateDump` if true, dump LFQ state on late packets (unexpected)
* `LateDumpPackets` if true, when dumping LFQ state on late packets, include
  info on each packet in the queues
* `FlowDefs` an array of flow definitions, each of which contains:
  * `Offset` offset from beginning of simulation to start enqueueing
  * `Interval` the average number of simulation ticks from one enqueue to the next
  * `IntervalVariance` the maximum number by which `Interval` can randomly vary
  * `Burst` the number of simultaneous packets to send at enqueue time
  * `BurstVariance` the maximum number by which `Burst` can randomly vary
  * `Size` the size of packets to send
  * `SizeVariance` the maximum number by which `Size` can randomly vary

## Results

The JSON output has the following format:

* `FlowStats` is an array of flow statistics, one for each flow in order, each of
  which contains:
  * `BytesSent` the total number of bytes sent for the flow
  * `Throughput` the flow's throughput, in bytes per 1000 ticks (thus, 0-1000)
  * `MeanSojourn` the mean sojourn time of all packets for the flow, in ticks
  * `MinSojourn` the minimum sojourn time of all packets for the flow, in ticks
  * `MaxSojourn` the maximum sojourn time of all packets for the flow, in ticks
  * `Enqueues` the number of packets that were passed to enqueue
  * `Drops` the number of packets that were dropped (`Enqueues` - `TotalSends`)
  * `DropsPercent` the percentage of enqueued packets that were dropped
  * `SparseSends` the number of packets sent from the sparse queue
  * `BulkSends` the number of packets sent from the bulk queue
  * `TotalSends` the total number of packets sent
  * `LateSends` the number of packets sent late (an out-of-order packets metric)
  * `LateSendsPercent` the percentage of packets sent late
