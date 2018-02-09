package algGen

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/MaxHalford/gago"
	uuid "github.com/nu7hatch/gouuid"
)

//NumFlags ...
const (
	NumFlags    = 141
	NoCorretude = 100000000
)

//Flags ...
type Flags []int

//Evaluate ...
func (f Flags) Evaluate() float64 {
	//log.Println("#### entrou no Evaluate")
	u, _ := uuid.NewV4()
	filename := "output/g" + u.String()
	//log.Println("#### antes do getTimeCompilation")
	_ = getTimeCompilation(f, filename)
	//log.Println("#### antes do getLengthCompilation")
	lenComp := getLengthCompilation(f, filename)
	//log.Println("#### depois do getLengthCompilation", lenCompilation)
	timeExec := getTimeExec(filename)
	if timeExec == NoCorretude {
		return NoCorretude
	}
	//log.Println("#### depois do getTimeExec", timeExec)
	return lenComp
}

func getTimeCompilation(f Flags, filename string) float64 {
	//cmd := exec.Command("/usr/bin/gccgo-6", "/home/raniere/Code/gccgo-test/main.go")
	//cmd := exec.Command("/usr/bin/gccgo-6", "alg/chain/expertseq/main.go")
	cmd := exec.Command("/usr/bin/gccgo-6", "alg/outer/expertseq/main.go")
	cmd.Env = os.Environ()

	for c, o := range f {
		if o == 1 {
			cmd.Args = append(cmd.Args, Opt[c])
		}
	}
	cmd.Args = append(cmd.Args, "-o")
	cmd.Args = append(cmd.Args, filename)

	start := time.Now()

	err := cmd.Start()
	if err != nil {
		log.Println("Erro Start:", err)
		return 0
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Erro Wait:", err)
		return 0
	}

	elapsed := time.Since(start)

	return float64(elapsed.Seconds())
}

func getLengthCompilation(f Flags, filename string) float64 {
	out, err := exec.Command("wc", "-c", filename).Output()
	if err != nil {
		log.Fatal(err)
	}

	lengthString := strings.Split(string(out), " ")[0]

	lengthInt, _ := strconv.Atoi(lengthString)

	return float64(lengthInt)
}

func getTimeExec(filename string) float64 {
	var elapsed time.Duration
	start := time.Now()

	cmd := exec.Command("bash", "-c", filename+" <  alg/outer/expertseq/main.in")

	out, err := cmd.Output()
	if err != nil {
		log.Println(err.Error(), filename, out)
		return NoCorretude
	}

	elapsed = time.Since(start)

	if string(out) == OutOuter {
		log.Println("eh igual", filename)
		return float64(elapsed.Seconds())
	}

	log.Println("eh diferente", len(string(out)), len(OutOuter), filename)
	log.Println(string(out))
	return NoCorretude
}

// Mutate a slice of Positions by permuting it's values.
func (f Flags) Mutate(rng *rand.Rand) {
	gago.MutPermuteInt(f, 3, rng)
	//gago.MutPermute(f, 3, rng)
}

// Crossover a slice of Positions with another by applying partially mapped
// crossover.
func (f Flags) Crossover(Y gago.Genome, rng *rand.Rand) {
	gago.CrossGNXInt(f, Y.(Flags), 2, rng)
	//return Flags(o1), Flags(o2)
	/*var o1, o2 = gago.CrossPMXInt(f, Y.(Flags), rng)
	return Flags(o1), Flags(o2)*/
	/*var o1, o2 = gago.CrossPMX(f, Y.(Flags), rng)
	return o1.(Flags), o2.(Flags)*/
}

// Clone a slice of Positions.
func (f Flags) Clone() gago.Genome {
	var ff = make(Flags, len(f))
	copy(ff, f)
	return ff
}

// MakeBoard creates a random slices of positions by generating random number
// permutations in [0, N_QUEENS).
func MakeBoard(rng *rand.Rand) gago.Genome {
	var flags = make(Flags, NumFlags)
	for i, flag := range rng.Perm(NumFlags) {
		if flag%2 == 0 {
			flags[i] = 0
		} else {
			flags[i] = 1
		}
		//flags[i] = int(flag%2 == 0)
	}
	return gago.Genome(flags)
}
