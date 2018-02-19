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
	"log"
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
		//log.Println(read)
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

	/*f, err := os.Create("cpuprofile.txt")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
	log.Println("depois do cpu profile")*/

	fi, err := os.Create("memprofile.txt")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}

	//for ii := 0; ii < 100000; ii++ {

		nelts = read_integer()
		points = make([]Point, nelts)

		if !*is_bench {
			read_vector_of_points(nelts)
		}

		matrix, vector := Outer(points[0:nelts], nelts)

		if !*is_bench {
			//fmt.Printf("%d %d\n", nelts, nelts)
			for _, row := range matrix {
				for _, elem := range row {
					if elem != 0.0 {

					}
				}
				//fmt.Printf("\n")
			}
			//fmt.Printf("\n")

			//fmt.Printf("%d\n", nelts)
			for i := 0; i < nelts; i++ {
				if vector[i] != 0.0 {

				}
			}
			//fmt.Printf("\n")
		}
	//}
	runtime.GC()
	//pprof.WriteHeapProfile(fi)
	//pprof.
	pprof.Lookup("heap").WriteTo(fi, 1)
	fi.Close()

	s := new(runtime.MemStats)
	runtime.ReadMemStats(s)
	fmt.Println("\n# runtime.MemStats")
	fmt.Println("# Alloc = %d", s.Alloc)
	fmt.Println("# TotalAlloc = %d", s.TotalAlloc)
	fmt.Println("# Sys = %d", s.Sys)
	fmt.Println("# Lookups = %d", s.Lookups)
	fmt.Println("# Mallocs = %d", s.Mallocs)
	fmt.Println("# Frees = %d", s.Frees)

	fmt.Println("# HeapAlloc = %d", s.HeapAlloc)
	fmt.Println("# HeapSys = %d", s.HeapSys)
	fmt.Println("# HeapIdle = %d", s.HeapIdle)
	fmt.Println("# HeapInuse = %d", s.HeapInuse)
	fmt.Println("# HeapReleased = %d", s.HeapReleased)
	fmt.Println("# HeapObjects = %d", s.HeapObjects)

	fmt.Println("# Stack = %d / %d", s.StackInuse, s.StackSys)
	fmt.Println("# MSpan = %d / %d", s.MSpanInuse, s.MSpanSys)
	fmt.Println("# MCache = %d / %d", s.MCacheInuse, s.MCacheSys)
	fmt.Println("# BuckHashSys = %d", s.BuckHashSys)
	fmt.Println("# GCSys = %d", s.GCSys)
	fmt.Println("# OtherSys = %d", s.OtherSys)

	fmt.Println("# NextGC = %d", s.NextGC)
	fmt.Println( "# LastGC = %d", s.LastGC)
	fmt.Println( "# PauseNs = %d", s.PauseNs)
	fmt.Println( "# PauseEnd = %d", s.PauseEnd)
	fmt.Println( "# NumGC = %d", s.NumGC)
	fmt.Println( "# NumForcedGC = %d", s.NumForcedGC)
	fmt.Println( "# GCCPUFraction = %v", s.GCCPUFraction)
	fmt.Println("# DebugGC = %v", s.DebugGC)
}
