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
    "QuickPull": true,
    "MaxSize": 48000,
    "FlowDefs": [
        {
            "Offset": 0,
            "Interval": 2950,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1200,
            "SizeVariance": 99
        },
        {
            "Offset": 0,
            "Interval": 100,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
            "Offset": 0,
            "Interval": 1000,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
            "Offset": 0,
            "Interval": 975,
            "IntervalVariance": 50,
            "Burst": 2,
            "BurstVariance": 2,
            "Size": 200,
            "SizeVariance": 50
        },
        {
            "Offset": 0,
            "Interval": 750,
            "IntervalVariance": 0,
            "Burst": 1,
            "BurstVariance": 0,
            "Size": 50,
            "SizeVariance": 100
        },
        {
            "Offset": 0,
            "Interval": 32000,
            "IntervalVariance": 0,
            "Burst": 32,
            "BurstVariance": 0,
            "Size": 1500,
            "SizeVariance": 0
        },
        {
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
            "BytesSent": 1037245,
            "Throughput": 103.7245,
            "MeanSojourn": 889.8030127462341,
            "SparseSends": 736,
            "BulkSends": 127,
            "TotalSends": 863,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 1038000,
            "Throughput": 103.8,
            "MeanSojourn": 2624.566473988439,
            "SparseSends": 13,
            "BulkSends": 679,
            "TotalSends": 692,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 1038000,
            "Throughput": 103.8,
            "MeanSojourn": 1231.21387283237,
            "SparseSends": 398,
            "BulkSends": 294,
            "TotalSends": 692,
            "LateSends": 1,
            "LateSendsPercent": 0.001445086705202312
        },
        {
            "BytesSent": 648176,
            "Throughput": 64.8176,
            "MeanSojourn": 1527.6317578004325,
            "SparseSends": 696,
            "BulkSends": 2541,
            "TotalSends": 3237,
            "LateSends": 7,
            "LateSendsPercent": 0.002162496138399753
        },
        {
            "BytesSent": 191261,
            "Throughput": 19.1261,
            "MeanSojourn": 1415.5861696380334,
            "SparseSends": 393,
            "BulkSends": 3309,
            "TotalSends": 3702,
            "LateSends": 136,
            "LateSendsPercent": 0.036736898973527825
        },
        {
            "BytesSent": 474000,
            "Throughput": 47.4,
            "MeanSojourn": 367.0886075949367,
            "SparseSends": 313,
            "BulkSends": 3,
            "TotalSends": 316,
            "LateSends": 0,
            "LateSendsPercent": 0
        },
        {
            "BytesSent": 49800,
            "Throughput": 4.98,
            "MeanSojourn": 613.9277108433735,
            "SparseSends": 497,
            "BulkSends": 1,
            "TotalSends": 498,
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
* `QuickPull` if true, use an experimental pull method that copies the queue
  tail to the scan position and shortens the tail by one
* `MaxSize` the maximum size of the queue, in bytes
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
