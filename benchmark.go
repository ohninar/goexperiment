package main

import (
	"log"
	"os"
	"os/exec"
	"time"
	"flag"
)

func main() {
	generations := flag.Int("generations", 1000, "generations=1000")
	pathBin := flag.String("bin", "target/randmat/expertpar/main-18", "bin='/path/bin'")
	pathInput := flag.String("input", "target/randmat/expertpar/main.in", "input='/path/input'")

	flag.Parse()

	var result float64

	for i := 0; i < *generations; i++ {
		result += runner(*pathBin, *pathInput)
	}

	log.Println(result / float64(*generations))
}

func runner(pathBin string, pathInput string) float64 {
	var elapsed time.Duration

	start := time.Now()
	cmd := exec.Command("bash", "-c", pathBin + " < " + pathInput)
	elapsed = time.Since(start)

	_, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	return elapsed.Seconds()
}
