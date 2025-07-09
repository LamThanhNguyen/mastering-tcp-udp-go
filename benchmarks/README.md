# TCP vs UDP Benchmark

## How to Run

First, start a server in one terminal:

```bash
go run benchmarks/server/server.go -proto=tcp
# or
go run benchmarks/server/server.go -proto=udp
```

Then, run the client in another terminal:

```bash
go run benchmarks/client/client.go -proto=tcp -n=1000 -size=128 -csv=tcp_rtt.csv
# or
go run benchmarks/client/client.go -proto=udp -n=1000 -size=128 -csv=udp_rtt.csv
```

- -n is the number of messages.
- -size is the message size in bytes.
- -csv is an optional file to save RTTs for analysis.

## Output
The client prints throughput, and detailed RTT analytics (min/avg/stddev/percentiles/max).

The CSV can be loaded in Excel, Google Sheets, or plotted with Python/matplotlib for graphs.