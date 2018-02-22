/*
 * randmat: random number generation
 *
 * input:
 *   nrows, ncols: the number of rows and columns
 *   s: the seed
 *
 * output:
 *   martix: a nrows x ncols integer matrix
 *
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

type ByteMatrix struct {
	Rows, Cols int
	array      []byte
}

func NewByteMatrix(r, c int) *ByteMatrix {
	return &ByteMatrix{r, c, make([]byte, r*c)}
}

func (m *ByteMatrix) Row(i int) []byte {
	return m.array[i*m.Cols : (i+1)*m.Cols]
}

const (
	LCG_A = 1664525
	LCG_C = 1013904223
)

var (
	is_bench = flag.Bool("is_bench", false, "")
)

func randmat(nrows, ncols int, s uint32) *ByteMatrix {
	matrix := NewByteMatrix(nrows, ncols)

	work := make(chan int)

	go func() {
		for i := 0; i < nrows; i++ {
			work <- i
		}
		close(work)
	}()

	done := make(chan bool)
	NP := runtime.GOMAXPROCS(0)

	for i := 0; i < NP; i++ {
		go func() {
			for i := range work {
				seed := s + uint32(i)
				row := matrix.Row(i)
				for j := range row {
					seed = LCG_A*seed + LCG_C
					row[j] = byte(seed%100) % 100
				}
			}
			done <- true
		}()
	}

	for i := 0; i < NP; i++ {
		<-done
	}

	return matrix
}

func main() {
	flag.Parse()

	var nrows, ncols int
	var seed uint32

	fmt.Scan(&nrows)
	fmt.Scan(&ncols)
	fmt.Scan(&seed)
	matrix := randmat(nrows, ncols, seed)

	if !*is_bench {
		for i := 0; i < nrows; i++ {
			row := matrix.Row(i)
			for j := range row {
				fmt.Printf("%d ", row[j])
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}

	if os.Getenv("GOEXP_DEBUG") == "1" {
		SaveMemProfile(os.Args[0])
	}
}

func SaveMemProfile(nameFile string) error {
	f, err := os.Create(nameFile + "-mem-profile.out")
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := os.Create(nameFile + "-mem-profile.txt")
	if err != nil {
		return err
	}
	defer fi.Close()

	runtime.GC()
	pprof.Lookup("heap").WriteTo(f, 1)

	s := new(runtime.MemStats)
	runtime.ReadMemStats(s)

	fmt.Fprintf(fi, "\n# runtime.MemStats\n")
	fmt.Fprintf(fi, "# Alloc = %d\n", s.Alloc)
	fmt.Fprintf(fi, "# TotalAlloc = %d\n", s.TotalAlloc)
	fmt.Fprintf(fi, "# Sys = %d\n", s.Sys)
	fmt.Fprintf(fi, "# Lookups = %d\n", s.Lookups)
	fmt.Fprintf(fi, "# Mallocs = %d\n", s.Mallocs)
	fmt.Fprintf(fi, "# Frees = %d\n\n", s.Frees)

	fmt.Fprintf(fi, "# HeapAlloc = %d\n", s.HeapAlloc)
	fmt.Fprintf(fi, "# HeapSys = %d\n", s.HeapSys)
	fmt.Fprintf(fi, "# HeapIdle = %d\n", s.HeapIdle)
	fmt.Fprintf(fi, "# HeapInuse = %d\n", s.HeapInuse)
	fmt.Fprintf(fi, "# HeapReleased = %d\n", s.HeapReleased)
	fmt.Fprintf(fi, "# HeapObjects = %d\n", s.HeapObjects)

	fmt.Fprintf(fi, "# Stack = %d / %d\n", s.StackInuse, s.StackSys)
	fmt.Fprintf(fi, "# MSpan = %d / %d\n", s.MSpanInuse, s.MSpanSys)
	fmt.Fprintf(fi, "# MCache = %d / %d\n", s.MCacheInuse, s.MCacheSys)
	fmt.Fprintf(fi, "# BuckHashSys = %d\n", s.BuckHashSys)
	fmt.Fprintf(fi, "# GCSys = %d\n", s.GCSys)
	fmt.Fprintf(fi, "# OtherSys = %d\n\n", s.OtherSys)

	fmt.Fprintf(fi, "# NextGC = %d\n", s.NextGC)
	fmt.Fprintf(fi, "# LastGC = %d\n", s.LastGC)
	fmt.Fprintf(fi, "# PauseNs = %d\n", s.PauseNs)
	fmt.Fprintf(fi, "# PauseEnd = %d\n", s.PauseEnd)
	fmt.Fprintf(fi, "# NumGC = %d\n", s.NumGC)
	fmt.Fprintf(fi, "# NumForcedGC = %d\n", s.NumForcedGC)
	fmt.Fprintf(fi, "# GCCPUFraction = %v\n", s.GCCPUFraction)
	fmt.Fprintf(fi, "# DebugGC = %v\n", s.DebugGC)

	return nil
}
