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
 */
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
	Rows, Cols int
	array      []byte
}

func WrapBytes(r, c int, bytes []byte) *ByteMatrix {
	return &ByteMatrix{r, c, bytes}
}

func NewByteMatrix(r, c int) *ByteMatrix {
	return &ByteMatrix{r, c, make([]byte, r*c)}
}

func (m *ByteMatrix) Row(i int) []byte {
	return m.array[i*m.Cols : (i+1)*m.Cols]
}

func (m *ByteMatrix) Bytes() []byte {
	return m.array[0 : m.Rows*m.Cols]
}

var is_bench = flag.Bool("is_bench", false, "")
var matrix []byte
var mask [][]bool
var points []int

type WinnowPoints struct {
	m *ByteMatrix
	e []int // indexes into the ByteMatrix 'm'
}

func (p *WinnowPoints) Len() int {
	return len(p.e)
}

func (p *WinnowPoints) Swap(i, j int) {
	p.e[i], p.e[j] = p.e[j], p.e[i]
}

func (p *WinnowPoints) Less(i, j int) bool {
	return ArrayLess(p.m.array, p.e[i], p.e[j])
}

func ArrayLess(array []byte, x, y int) bool {
	if array[x] != array[y] {
		return array[x] < array[y]
	}
	return x < y
}

func WinnowMerge(points chan WinnowPoints) {
	var merged WinnowPoints
	x := <-points
	y := <-points

	new_size := len(x.e) + len(y.e)

	merged.m = x.m
	merged.e = make([]int, new_size)

	j := 0
	k := 0
	for i := 0; i < new_size; i++ {
		if j < len(x.e) && k < len(y.e) {
			if ArrayLess(merged.m.array, x.e[j], y.e[k]) {
				merged.e[i] = x.e[j]
				j++
			} else {
				merged.e[i] = y.e[k]
				k++
			}
		} else if j < len(x.e) {
			merged.e[i] = x.e[j]
			j++
		} else if k < len(y.e) {
			merged.e[i] = y.e[k]
			k++
		}
	}
	points <- merged
}

func Winnow(m *ByteMatrix, nrows, ncols, winnow_nelts int) {
	NP := runtime.GOMAXPROCS(0)
	var values WinnowPoints
	values.m = m

	values_work := make(chan int, 1024)
	values_done := make(chan WinnowPoints, NP)
	values_done <- WinnowPoints{m, make([]int, 0)}

	go func() {
		for i := 0; i < nrows; i++ {
			values_work <- i
		}
		close(values_work)
	}()

	merged := make(chan bool, NP)

	for i := 0; i < NP; i++ {
		go func() {
			WinnowMerge(values_done)
			merged <- true
		}()
	}

	for i := 0; i < NP; i++ {
		go func() {
			var local_indexes []int
			for i := range values_work {
				for j := 0; j < ncols; j++ {
					idx := i*ncols + j
					if *is_bench {
						mask[i][j] = ((i * j) % (ncols + 1)) == 1
					}
					if mask[i][j] {
						local_indexes = append(local_indexes, idx)
					}
				}
			}
			var local_values WinnowPoints
			local_values.m = m
			local_values.e = local_indexes

			sort.Sort(&local_values)
			values_done <- local_values
		}()
	}

	for i := 0; i < NP; i++ {
		<-merged
	}

	values = <-values_done

	chunk := values.Len() / winnow_nelts

	points = make([]int, winnow_nelts)
	point_work := make(chan int, 1024)
	point_done := make(chan bool)
	go func() {
		for i := 0; i < winnow_nelts; i++ {
			point_work <- i
		}
		close(point_work)
	}()

	for i := 0; i < NP; i++ {
		go func() {
			for i := range point_work {
				points[i] = values.e[i*chunk]
			}
			point_done <- true
		}()
	}

	for i := 0; i < NP; i++ {
		<-point_done
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

func read_matrix(nrows, ncols int) {
	for i := 0; i < nrows; i++ {
		for j := 0; j < ncols; j++ {
			matrix[i*ncols+j] = byte(read_integer())
		}
	}
}

func read_mask(nrows, ncols int) {
	for i := 0; i < nrows; i++ {
		for j := 0; j < ncols; j++ {
			mask[i][j] = (read_integer() == 1)
		}
	}
}

func main() {
	var nrows, ncols, nelts int

	flag.Parse()

	nrows = int(read_integer())
	ncols = int(read_integer())

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

	nelts = int(read_integer())
	points = make([]int, nelts)

	Winnow(m, nrows, ncols, nelts)

	if !*is_bench {
		fmt.Printf("%d\n", nelts)
		for i := 0; i < nelts; i++ {
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
