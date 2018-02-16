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

func Outer(wp []Point, nelts int) (m []float64, vec []float64) {
	m = make([]float64, nelts*nelts)
	vec = make([]float64, nelts)
	for i, v := range wp {
		nmax := float64(0)
		for j, w := range wp {
			if i != j {
				d := Distance(v.x, v.y, w.x, w.y)
				if d > nmax {
					nmax = d
				}
				m[i*nelts+j] = d
			}
		}
		m[i*(nelts+1)] = float64(nelts) * nmax
		vec[i] = Distance(0, 0, v.x, v.y)
	}
	return
}

func read_integer() int {
	var value int
	for true {
		var read, _ = fmt.Scanf("%d", &value)
		log.Println("read:", read)
		if read == 1 || read == 0 {
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

	f, err := os.Create("cpuprofile.txt")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	nelts = read_integer()
	points = make([]Point, nelts)

	log.Println("nelts:", nelts)

	if !*is_bench {
		read_vector_of_points(nelts)
	}

	matrix, vector := Outer(points[0:nelts], nelts)

	for ii := 0; ii < 1000; ii++ {

		if !*is_bench {
			fmt.Printf("%d %d", nelts, nelts)
			for i := 0; i < nelts*nelts; i++ {
				if i%nelts == 0 {
					fmt.Printf("\n")
				}
				fmt.Printf("%g ", matrix[i])
			}
			fmt.Printf("\n\n")

			fmt.Printf("%d\n", nelts)
			for i := 0; i < nelts; i++ {
				fmt.Printf("%g ", vector[i])
			}
			fmt.Printf("\n")
		}
	}

	fi, err := os.Create("memprofile.txt")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(fi); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	fi.Close()
}
