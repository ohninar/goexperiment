package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/montanaflynn/stats"
)

func main() {
	generations := flag.Int("generations", 1000, "generations=1000")
	pathBin := flag.String("bin", "", "bin='/path/bin'")
	pathInput := flag.String("input", "", "input='/path/input'")

	flag.Parse()

	fmt.Println(*generations, *pathBin, *pathInput)

	var elapsed time.Duration
	var results []float64

	start := time.Now()
	fmt.Println("Inicio:", start)

	for i := 0; i < *generations; i++ {
		results = append(results, runner2(*pathBin, *pathInput))
	}

	elapsed = time.Since(start)
	theEnd := time.Now()
	fmt.Println("   Fim:", theEnd)
	fmt.Println("  Diff:", elapsed)

	median, _ := stats.Median(results)
	fmt.Printf("\n Mediana: %.10f", median)
}

func runner2(pathBin string, pathInput string) float64 {
	var elapsed time.Duration

	command := pathBin + " < " + pathInput
	if pathInput == "" {
		command = pathBin
	}

	start := time.Now()
	_, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	elapsed = time.Since(start)

	return elapsed.Seconds()
}

func runner(pathBin string, pathInput string) float64 {
	var elapsed time.Duration

	command := pathBin + " < " + pathInput
	if pathInput == "" {
		command = pathBin
	}

	start := time.Now()
	cmd := exec.Command("bash", "-c", command)
	elapsed = time.Since(start)

	_, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	return elapsed.Seconds()
}
