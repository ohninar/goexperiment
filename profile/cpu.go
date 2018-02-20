package profile

import (
	"os"
	"runtime/pprof"
)

func SaveCPUProfile(nameFile string) error {
	f, err := os.Create(nameFile + "-cpu-profile.out")
	if err != nil {
		return err
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		return err
	}

	return nil
}