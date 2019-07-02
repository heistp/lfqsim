# lfqsim

This is a discrete time simulator for Lightweight Fair Queueing.

## Quick Start

`lfqsim` reads its configuration from stdin and writes results to stdout, both in
JSON format. After having installed Go using your package manager or from
[here](https://golang.org/dl/):

```
go get github.com/heistp/lfqsim
cd ~/go/src/github.com/heistp/lfqsim # or location of lfqsim directory
go build
./lfqsim < config.json
```

## Description

Lightweight Fair Queueing is a fair queueing algorithm from Jonathon Morton with
a small codebase, low memory requirements, and a design suitable for
implementation in hardware. This related project provides an LFQ implementation
(in `lfq.go`) and a discrete time simulation that measures how it performs with
a configurable number of flows and their characteristics.

Each unit of time is referred to as a tick. Enqueue is called for each
configured flow at regular or varying intervals in ticks. Dequeue is called at a
constant bitrate of one tick per byte. See the [Results](#results) section for
information on what results are produced. Note that late packets (out of order
packets) are not expected with LFQ, so please file an issue if you discover
simulation parameters that cause them.

## Differences from Internet Draft Specification

- No AQM (Active Queue Management), so we rely on overflow to manage queue size
- The packet hash is specified directly
- For simulation purposes:
  - Timestamp is a Tick
  - Send method contains sparse flag for stats

## Sample Simulation and Discussion

Using the default `config.json`, output similar to what's below may be produced.

```
{
    "EndTicks": 100000000,
    "MTU": 1500,
    "MaxSize": 256000,
    "LateDump": true,
    "LateDumpPackets": false,
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
            "BytesSent": 22463044,
            "Throughput": 224.63044,
            "MeanSojourn": 46936.56279044923,
            "MinSojourn": 7570,
            "MaxSojourn": 50811,
            "Enqueues": 27028,
            "Drops": 8307,
            "DropsPercent": 30.734793547432293,
            "SparseSends": 1,
            "BulkSends": 18720,
            "TotalSends": 18721,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 22464000,
            "Throughput": 224.64,
            "MeanSojourn": 48914.4421073718,
            "MinSojourn": 1278,
            "MaxSojourn": 50904,
            "Enqueues": 200000,
            "Drops": 185024,
            "DropsPercent": 92.512,
            "SparseSends": 1,
            "BulkSends": 14975,
            "TotalSends": 14976,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 22464000,
            "Throughput": 224.64,
            "MeanSojourn": 48675.2687633547,
            "MinSojourn": 2778,
            "MaxSojourn": 50999,
            "Enqueues": 66667,
            "Drops": 51691,
            "DropsPercent": 77.5361123194384,
            "SparseSends": 1,
            "BulkSends": 14975,
            "TotalSends": 14976,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 22462500,
            "Throughput": 224.625,
            "MeanSojourn": 24429.68567612688,
            "MinSojourn": 3,
            "MaxSojourn": 52862,
            "Enqueues": 66688,
            "Drops": 51713,
            "DropsPercent": 77.54468570057581,
            "SparseSends": 1,
            "BulkSends": 14974,
            "TotalSends": 14975,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 2981035,
            "Throughput": 29.81035,
            "MeanSojourn": 3910.195986307806,
            "MinSojourn": 2,
            "MaxSojourn": 14312,
            "Enqueues": 14901,
            "Drops": 2,
            "DropsPercent": 0.013421917992081069,
            "SparseSends": 6429,
            "BulkSends": 8470,
            "TotalSends": 14899,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 6665988,
            "Throughput": 66.65988,
            "MeanSojourn": 5477.608669236885,
            "MinSojourn": 1,
            "MaxSojourn": 13287,
            "Enqueues": 133334,
            "Drops": 12,
            "DropsPercent": 0.008999955000225,
            "SparseSends": 8050,
            "BulkSends": 125272,
            "TotalSends": 133322,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 500000,
            "Throughput": 5,
            "MeanSojourn": 663.662,
            "MinSojourn": 1,
            "MaxSojourn": 6065,
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
the results. From the results, it can be seen that throughput fairness is
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
