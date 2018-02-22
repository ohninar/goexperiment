/*
 * winnow: weighted point selection
 *
 * input:
 *   matrix: an integer matrix, whose values are used as masses
 *   mask: a boolean matrix showing which points are eligible for
 *     consideration
 *   nrows, ncols: the number of rows and columns
 *   nelts: the number of points to select
 *
 * output:
 *   points: a vector of (x, y) points
 *
 :b*/
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
)

type ByteMatrix struct {
	Rows, Cols uint32
	array      []byte
}

func WrapBytes(r, c uint32, bytes []byte) *ByteMatrix {
	return &ByteMatrix{r, c, bytes}
}

func NewByteMatrix(r, c uint32) *ByteMatrix {
	return &ByteMatrix{r, c, make([]byte, r*c)}
}

func (m *ByteMatrix) Row(i uint32) []byte {
	return m.array[i*m.Cols : (i+1)*m.Cols]
}

func (m *ByteMatrix) Bytes() []byte {
	return m.array[0 : m.Rows*m.Cols]
}

var is_bench = flag.Bool("is_bench", false, "")
var matrix []byte
var mask [][]bool
var points []uint32

type WinnowPoints struct {
	m *ByteMatrix
	e []uint32 // indexes into the ByteMatrix 'm'
}

func (p *WinnowPoints) Len() int {
	return len(p.e)
}

func (p *WinnowPoints) Swap(i, j int) {
	p.e[i], p.e[j] = p.e[j], p.e[i]
}

func (p *WinnowPoints) Less(i, j int) bool {
	if p.m.array[p.e[i]] != p.m.array[p.e[j]] {
		return p.m.array[p.e[i]] < p.m.array[p.e[j]]
	}

	return p.e[i] < p.e[j]
}

func Winnow(m *ByteMatrix, nrows, ncols, nelts uint32) {
	var values WinnowPoints
	values.m = m

	for i := uint32(0); i < nrows; i++ {
		for j := uint32(0); j < ncols; j++ {
			if *is_bench {
				mask[i][j] = ((i * j) % (ncols + 1)) == 1
			}

			if mask[i][j] {
				idx := i*ncols + j
				values.e = append(values.e, idx)
			}
		}
	}

	sort.Sort(&values)

	chunk := uint32(values.Len()) / nelts

	for i := uint32(0); i < nelts; i++ {
		points[i] = values.e[i*chunk]
	}
}

func read_integer() int {
	var value int
	for true {
		var read, _ = fmt.Scanf("%d", &value)
		if read == 1 {
			break
		}
	}
	return value
}

func read_matrix(nrows, ncols uint32) {
	for i := uint32(0); i < nrows; i++ {
		for j := uint32(0); j < ncols; j++ {
			matrix[i*ncols+j] = byte(read_integer())
		}
	}
}

func read_mask(nrows, ncols uint32) {
	for i := uint32(0); i < nrows; i++ {
		for j := uint32(0); j < ncols; j++ {
			mask[i][j] = (read_integer() == 1)
		}
	}
}

func main() {
	var nrows, ncols, nelts uint32

	flag.Parse()

	nrows = uint32(read_integer())
	ncols = uint32(read_integer())

	m := NewByteMatrix(nrows, ncols)

	matrix = m.array

	mask = make([][]bool, nrows)
	for i := range mask {
		mask[i] = make([]bool, ncols)
	}

	if !*is_bench {
		read_matrix(nrows, ncols)
		read_mask(nrows, ncols)
	}

	nelts = uint32(read_integer())
	points = make([]uint32, nelts)

	Winnow(m, nrows, ncols, nelts)

	if !*is_bench {
		fmt.Printf("%d\n", nelts)
		for i := uint32(0); i < nelts; i++ {
			fmt.Printf("%d %d\n", points[i]/ncols, points[i]%ncols)
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
