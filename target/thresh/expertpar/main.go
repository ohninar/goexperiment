/*
 * thresh: histogram thresholding
 *
 * input:
 *   matrix: the integer matrix to be thresholded
 *   nrows, ncols: the number of rows and columns
 *   percent: the percentage of cells to retain
 *
 * output:
 *   mask: a boolean matrix filled with true for cells that are kept
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

var is_bench = flag.Bool("is_bench", false, "")

type ByteMatrix struct {
	Rows, Cols uint32
	array      []byte
}

func NewByteMatrix(r, c uint32) *ByteMatrix {
	return &ByteMatrix{r, c, make([]byte, r*c)}
}

func WrapBytes(r, c uint32, bytes []byte) *ByteMatrix {
	return &ByteMatrix{r, c, bytes}
}

func (m *ByteMatrix) Row(i uint32) []byte {
	return m.array[i*m.Cols : (i+1)*m.Cols]
}

func (m *ByteMatrix) Bytes() []byte {
	return m.array[0 : m.Rows*m.Cols]
}

var mask [][]bool

func thresh(m *ByteMatrix, nrows, ncols, percent uint32) {
	NP := runtime.GOMAXPROCS(0)

	hist_work := make(chan uint32)
	hist_parts := make(chan []int)
	go func() {
		for i := uint32(0); i < nrows; i++ {
			hist_work <- i
		}
		close(hist_work)
	}()

	for i := 0; i < NP; i++ {
		go func() {
			my_hist := make([]int, 100)
			for i := range hist_work {
				row := m.Row(i)
				for j := range row {
					my_hist[row[j]]++
				}
			}
			hist_parts <- my_hist
		}()
	}

	var hist [100]int

	for i := 0; i < NP; i++ {
		my_hist := <-hist_parts
		for j := range my_hist {
			hist[j] += my_hist[j]
		}
	}

	count := (nrows * ncols * percent) / 100
	prefixsum := 0
	threshold := 99

	for ; threshold > 0; threshold-- {
		prefixsum += hist[threshold]
		if prefixsum > int(count) {
			break
		}
	}

	mask_work := make(chan uint32)

	go func() {
		for i := uint32(0); i < nrows; i++ {
			mask_work <- i
		}
		close(mask_work)
	}()

	mask_done := make(chan bool)
	for i := 0; i < NP; i++ {
		go func() {
			for i := range mask_work {
				row := m.Row(i)
				for j := range row {
					mask[i][j] = row[j] >= byte(threshold)
				}
			}
			mask_done <- true
		}()
	}

	for i := 0; i < NP; i++ {
		<-mask_done
	}
}

func main() {
	var nrows, ncols, percent uint32

	flag.Parse()

	fmt.Scanf("%d%d", &nrows, &ncols)
	mask = make([][]bool, nrows)
	for i := range mask {
		mask[i] = make([]bool, ncols)
	}

	m := WrapBytes(nrows, ncols, make([]byte, ncols*nrows))

	if !*is_bench {
		for i := uint32(0); i < nrows; i++ {
			row := m.Row(i)
			for j := range row {
				fmt.Scanf("%d", &row[j])
			}
		}
	}

	fmt.Scanf("%d", &percent)

	thresh(m, nrows, ncols, percent)

	if !*is_bench {
		for i := uint32(0); i < nrows; i++ {
			for j := uint32(0); j < ncols; j++ {
				if mask[i][j] {
					fmt.Printf("1 ")
				} else {
					fmt.Printf("0 ")
				}
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
