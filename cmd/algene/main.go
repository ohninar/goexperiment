package main

import (
"log"
"os"

"github.com/MaxHalford/gago"
"github.com/raninho/goexperiment/algGen"
)

func main() {
	f, err := os.OpenFile("testlogfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("erro ao abrir arquivo do log.")
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("Iniciando o Initialize")
	var ga = gago.Generational(algGen.MakeBoard)
	ga.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	ga.Initialize()
	log.Println("Finalizando o Initialize")

	const GENERATIONS = 3000

	for i := 0; i < GENERATIONS; i++ {
		ga.Evolve()
		log.Println("Finalizando o Enhance com fitness", i, ga.HallOfFame[0].Fitness, ga.HallOfFame[0].Genome)
	}

	log.Printf("Optimal solution obtained after %d generations in %s\n", ga.Generations, ga.Age)
}
