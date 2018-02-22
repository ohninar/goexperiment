/*
 * outer: outer product
 *
 * input:
 *   vector: a vector of (x, y) points
 *   nelts: the number of points
 *
 * output:
 *   matrix: a real matrix, whose values are filled with inter-point
 *     distances
 *   vector: a real vector, whose values are filled with origin-to-point
 *     distances
 */
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
)

var is_bench = flag.Bool("is_bench", false, "")

var points []Point

type Point struct {
	x, y int
}

func Sqr(x float64) float64 {
	return x * x
}

func Distance(ax, ay, bx, by int) float64 {
	return math.Sqrt(float64(Sqr(float64(ax-bx)) + Sqr(float64(ay-by))))
}

func Outer(wp []Point, nelts int) (m [][]float64, vec []float64) {
	m = make([][]float64, nelts)
	vec = make([]float64, nelts)

	NP := runtime.GOMAXPROCS(0)
	work := make(chan int)

	done := make(chan bool)

	go func() {
		for i := range wp {
			work <- i
		}
		close(work)
	}()

	for i := 0; i < NP; i++ {
		go func() {
			for i := range work {
				m[i] = make([]float64, nelts)
				v := wp[i]
				nmax := float64(0)
				for j, w := range wp {
					if i != j {
						d := Distance(v.x, v.y, w.x, w.y)
						if d > nmax {
							nmax = d
						}
						m[i][j] = d
					}
				}
				m[i][i] = float64(nelts) * nmax
				vec[i] = Distance(0, 0, v.x, v.y)
			}
			done <- true
		}()
	}

	for i := 0; i < NP; i++ {
		<-done
	}
	return
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

func read_vector_of_points(nelts int) {
	for i := 0; i < nelts; i++ {
		a := read_integer()
		b := read_integer()
		points[i] = Point{a, b}
	}
}

func main() {
	var nelts int

	flag.Parse()

	nelts = read_integer()
	points = make([]Point, nelts)

	if !*is_bench {
		read_vector_of_points(nelts)
	}

	matrix, vector := Outer(points[0:nelts], nelts)

	if !*is_bench {
		fmt.Printf("%d %d\n", nelts, nelts)
		for _, row := range matrix {
			for _, elem := range row {
				fmt.Printf("%g ", elem)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")

		fmt.Printf("%d\n", nelts)
		for i := 0; i < nelts; i++ {
			fmt.Printf("%g ", vector[i])
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
