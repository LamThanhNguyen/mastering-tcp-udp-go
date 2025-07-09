package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	tcpAddr = "127.0.0.1:9101"
	udpAddr = "127.0.0.1:9102"
)

func main() {
	proto := flag.String("proto", "tcp", "tcp or udp")
	n := flag.Int("n", 1000, "number of messages to send")
	size := flag.Int("size", 64, "message size in bytes")
	csvFile := flag.String("csv", "", "CSV file to save per-message RTTs (optional)")
	flag.Parse()

	if *proto == "tcp" {
		runTCPClientRTT(*n, *size, *csvFile)
	} else {
		runUDPClientRTT(*n, *size, *csvFile)
	}
}

// ---------------- TCP RTT ------------------

func runTCPClientRTT(n, size int, csvFile string) {
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	msg := strings.Repeat("A", size)
	rtts := make([]float64, n)
	start := time.Now()
	for i := range n {
		sendTime := time.Now()
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Write error:", err)
			continue
		}
		buf := make([]byte, size)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		rtts[i] = float64(time.Since(sendTime).Microseconds())
	}
	elapsed := time.Since(start)
	printStats("TCP", n, size, elapsed, rtts)
	if csvFile != "" {
		writeCSV(csvFile, rtts)
	}
}

// ---------------- UDP RTT ------------------

func runUDPClientRTT(n, size int, csvFile string) {
	addr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	msg := []byte(strings.Repeat("B", size))
	rtts := make([]float64, n)
	start := time.Now()
	for i := range n {
		sendTime := time.Now()
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("Write error:", err)
			continue
		}
		buf := make([]byte, size)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		rtts[i] = float64(time.Since(sendTime).Microseconds())
	}
	elapsed := time.Since(start)
	printStats("UDP", n, size, elapsed, rtts)
	if csvFile != "" {
		writeCSV(csvFile, rtts)
	}
}

// -------------- Stats, CSV, Analytics --------------

func printStats(proto string, n, size int, elapsed time.Duration, rtts []float64) {
	validRTTs := rtts[:0]
	for _, v := range rtts {
		if v > 0 {
			validRTTs = append(validRTTs, v)
		}
	}
	if len(validRTTs) == 0 {
		fmt.Println("No RTTs recorded.")
		return
	}
	sort.Float64s(validRTTs)
	min := validRTTs[0]
	max := validRTTs[len(validRTTs)-1]
	sum := 0.0
	for _, v := range validRTTs {
		sum += v
	}
	avg := sum / float64(len(validRTTs))
	std := stddev(validRTTs, avg)
	p50 := percentile(validRTTs, 50)
	p95 := percentile(validRTTs, 95)
	p99 := percentile(validRTTs, 99)

	fmt.Printf("%s RTT Benchmark:\n", proto)
	fmt.Printf("Sent %d messages of %d bytes in %v\n", n, size, elapsed)
	fmt.Printf("Throughput: %.2f msg/sec, %.2f MB/sec\n",
		float64(n)/elapsed.Seconds(), float64(n*size)/(1024*1024)/elapsed.Seconds())
	fmt.Printf("RTT (min/avg/stddev/p50/p95/p99/max): %.2f / %.2f / %.2f / %.2f / %.2f / %.2f / %.2f us\n",
		min, avg, std, p50, p95, p99, max)
}

func stddev(vals []float64, avg float64) float64 {
	sum := 0.0
	for _, v := range vals {
		sum += (v - avg) * (v - avg)
	}
	return math.Sqrt(sum / float64(len(vals)))
}

func percentile(vals []float64, p float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	idx := max(int(math.Ceil((p/100)*float64(len(vals))))-1, 0)
	return vals[idx]
}

func writeCSV(filename string, vals []float64) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create CSV file: %v\n", err)
		return
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	w.Write([]string{"RTT_us"})
	for _, v := range vals {
		w.Write([]string{fmt.Sprintf("%.2f", v)})
	}
	fmt.Printf("RTTs written to %s\n", filename)
}
