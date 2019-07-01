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
            "BytesSent": 1005473,
            "Throughput": 100.5473,
            "MeanSojourn": 922.3747016706444,
            "Enqueues": 3390,
            "Drops": 2552,
            "DropsPercent": 75.28023598820059,
            "SparseSends": 699,
            "BulkSends": 139,
            "TotalSends": 838,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 1006500,
            "Throughput": 100.65,
            "MeanSojourn": 2599.2548435171384,
            "Enqueues": 100000,
            "Drops": 99329,
            "DropsPercent": 99.329,
            "SparseSends": 8,
            "BulkSends": 663,
            "TotalSends": 671,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 1006500,
            "Throughput": 100.65,
            "MeanSojourn": 2025.3353204172877,
            "Enqueues": 10000,
            "Drops": 9329,
            "DropsPercent": 93.29,
            "SparseSends": 116,
            "BulkSends": 555,
            "TotalSends": 671,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 498000,
            "Throughput": 49.8,
            "MeanSojourn": 240.96385542168676,
            "Enqueues": 10016,
            "Drops": 9684,
            "DropsPercent": 96.685303514377,
            "SparseSends": 313,
            "BulkSends": 19,
            "TotalSends": 332,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 627792,
            "Throughput": 62.7792,
            "MeanSojourn": 1440.6107746254384,
            "Enqueues": 15429,
            "Drops": 12292,
            "DropsPercent": 79.6681573659991,
            "SparseSends": 776,
            "BulkSends": 2361,
            "TotalSends": 3137,
            "LateSends": 9,
            "LateSendsPercent": 0.2868983104877271
        },
        {
            "BytesSent": 192191,
            "Throughput": 19.2191,
            "MeanSojourn": 1241.1177385892115,
            "Enqueues": 13334,
            "Drops": 9478,
            "DropsPercent": 71.08144592770361,
            "SparseSends": 356,
            "BulkSends": 3500,
            "TotalSends": 3856,
            "LateSends": 28,
            "LateSendsPercent": 0.7261410788381742
        },
        {
            "BytesSent": 49500,
            "Throughput": 4.95,
            "MeanSojourn": 740.9151515151515,
            "Enqueues": 500,
            "Drops": 5,
            "DropsPercent": 1,
            "SparseSends": 495,
            "BulkSends": 0,
            "TotalSends": 495,
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
  * `SparseSends` the number of packets sent from the sparse queue
  * `BulkSends` the number of packets sent from the bulk queue
  * `TotalSends` the total number of packets sent
  * `LateSends` the number of packets sent late (an out-of-order packets metric)
  * `LateSendsPercent` the percentage of packets sent late
