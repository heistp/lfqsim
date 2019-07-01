# lfqsim

This is a discrete time simulator for Lightweight Fair Queueing. It takes its
configuration via JSON to stdin, and writes results using JSON to stdout.

## Quick Start

After having installed Go:

```
go get https://github.com/heistp/lfqsim
cd lfqsim # change to location of lfqsim directory
go build
./lfqsim < config.json
```

## Sample Simulation

Using the default config.json, output similar to the following may be produced:

```
{
    "EndTicks": 10000000,
    "DequeueInterval": 1000,
    "MTU": 1500,
    "FastPull": false,
    "MaxSize": 48000,
    "LateDump": true,
    "LateDumpPackets": false,
    "FlowDefs": [
        {
            "Description": "1200 byte packets at just less then 3x dequeue interval",
            "Offset": 0,
            "Interval": 2950,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1200,
            "SizeVariance": 99
        },
        {
            "Description": "Aggressive flow- 1500 byte packets at 1/10 dequeue interval",
            "Offset": 0,
            "Interval": 100,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
            "Description": "1500 byte packets exactly at dequeue interval",
            "Offset": 0,
            "Interval": 1000,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
            "Description": "Bursty flow- 1500 byte packets with 32 packet bursts (exactly MaxSize) at 32x dequeue interval",
            "Offset": 0,
            "Interval": 32000,
            "IntervalVariance": 0,
            "Burst": 32,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
            "Description": "Smallish packets- 200 byte packets, burst 2 +/- 2, at just short of dequeue interval (causes re-ordering)",
            "Offset": 0,
            "Interval": 975,
            "IntervalVariance": 50,
            "Burst": 2,
            "BurstVariance": 2,
            "Size": 200,
            "SizeVariance": 50
        },
        {
            "Description": "Small packets- 50 byte packets at 3/4 dequeue interval (causes re-ordering)",
            "Offset": 0,
            "Interval": 750,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 50,
            "SizeVariance": 10
        },
        {
            "Description": "Sparse flow- 100 byte packets at 20x dequeue interval",
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
            "BytesSent": 442545,
            "Throughput": 44.2545,
            "MeanSojourn": 1386.2771739130435,
            "Enqueues": 3390,
            "Drops": 3022,
            "DropsPercent": 89.14454277286136,
            "SparseSends": 306,
            "BulkSends": 62,
            "TotalSends": 368,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 444000,
            "Throughput": 44.4,
            "MeanSojourn": 2186.1486486486488,
            "Enqueues": 100000,
            "Drops": 99704,
            "DropsPercent": 99.704,
            "SparseSends": 81,
            "BulkSends": 215,
            "TotalSends": 296,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 444000,
            "Throughput": 44.4,
            "MeanSojourn": 2344.5945945945946,
            "Enqueues": 10000,
            "Drops": 9704,
            "DropsPercent": 97.04,
            "SparseSends": 168,
            "BulkSends": 128,
            "TotalSends": 296,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 444000,
            "Throughput": 44.4,
            "MeanSojourn": 138.51351351351352,
            "Enqueues": 10016,
            "Drops": 9720,
            "DropsPercent": 97.04472843450479,
            "SparseSends": 213,
            "BulkSends": 83,
            "TotalSends": 296,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 425491,
            "Throughput": 42.5491,
            "MeanSojourn": 2074.2936470588234,
            "Enqueues": 15561,
            "Drops": 13436,
            "DropsPercent": 86.34406529143371,
            "SparseSends": 469,
            "BulkSends": 1656,
            "TotalSends": 2125,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 312175,
            "Throughput": 31.2175,
            "MeanSojourn": 2214.5319012504006,
            "Enqueues": 13334,
            "Drops": 7096,
            "DropsPercent": 53.217339133043346,
            "SparseSends": 377,
            "BulkSends": 5861,
            "TotalSends": 6238,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 38100,
            "Throughput": 3.81,
            "MeanSojourn": 783.8818897637796,
            "Enqueues": 501,
            "Drops": 120,
            "DropsPercent": 23.952095808383234,
            "SparseSends": 375,
            "BulkSends": 6,
            "TotalSends": 381,
            "LateSends": 0,
            "LateSendsPercent": 0
        }
    ]
}
```

## Configuration

The `config.json` file controls the simulation, and has the following format:

* `DequeueInterval` the number of simulation ticks from one dequeue to the next
* `EndTicks` the last tick for the simulation, exclusive
* `MTU` the MTU used for LFQ
* `FastPull` if true, use an experimental pull method that copies the queue
  tail to the scan position and shortens the tail by one
* `MaxSize` the maximum size of the queue, in bytes
* `LateDump` if true, dump LFQ state on late packets
* `LateDumpPackets` if true, when dumping LFQ state on late packets, include
  info on each packet in the queues
* `FlowDefs` an array of flow definitions, each of which contains:
  * `Offset` offset from beginning of simulation to start enqueueing
  * `Interval` the number of simulation ticks from one enqueue to the next
  * `IntervalVariance` the maximum number by which `Interval` can randomly vary
  * `Burst` the number of simultaneous packets to send at enqueue time
  * `BurstVariance` the maximum number by which `Burst` can randomly vary
  * `Size` the size of packets to send
  * `SizeVariance` the maximum number by which `Size` can randomly vary

## Results

The JSON output has the following format:

* `FlowStats` an array of flow statistics, one for each flow in order, each of
  which contains:
  * `BytesSent` the total number of bytes sent for the flow
  * `Throughput` the flow's throughput, in bytes per dequeue interval
  * `MeanSojourn` the mean sojourn time of all packets, in ticks
  * `Enqueues` the number of packets that were passed to enqueue
  * `Drops` the number of packets that were dropped (`Enqueues` - `TotalSends`)
  * `DropsPercent` the percentage of enqueued packets that were dropped
  * `SparseSends` the number of packets sent from the sparse queue
  * `BulkSends` the number of packets sent from the bulk queue
  * `TotalSends` the total number of packets sent
  * `LateSends` the number of packets sent late (an out-of-order packets metric)
  * `LateSendsPercent` the percentage of packets sent late
