package profile

import (
	"os"
	"runtime"
	"runtime/pprof"
	"fmt"
)

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

	fmt.Fprintf(fi,"\n# runtime.MemStats\n")
	fmt.Fprintf(fi,"# Alloc = %d\n", s.Alloc)
	fmt.Fprintf(fi,"# TotalAlloc = %d\n", s.TotalAlloc)
	fmt.Fprintf(fi,"# Sys = %d\n", s.Sys)
	fmt.Fprintf(fi,"# Lookups = %d\n", s.Lookups)
	fmt.Fprintf(fi,"# Mallocs = %d\n", s.Mallocs)
	fmt.Fprintf(fi,"# Frees = %d\n\n", s.Frees)

	fmt.Fprintf(fi,"# HeapAlloc = %d\n", s.HeapAlloc)
	fmt.Fprintf(fi,"# HeapSys = %d\n", s.HeapSys)
	fmt.Fprintf(fi,"# HeapIdle = %d\n", s.HeapIdle)
	fmt.Fprintf(fi,"# HeapInuse = %d\n", s.HeapInuse)
	fmt.Fprintf(fi,"# HeapReleased = %d\n", s.HeapReleased)
	fmt.Fprintf(fi,"# HeapObjects = %d\n", s.HeapObjects)

	fmt.Fprintf(fi,"# Stack = %d / %d\n", s.StackInuse, s.StackSys)
	fmt.Fprintf(fi,"# MSpan = %d / %d\n", s.MSpanInuse, s.MSpanSys)
	fmt.Fprintf(fi,"# MCache = %d / %d\n", s.MCacheInuse, s.MCacheSys)
	fmt.Fprintf(fi,"# BuckHashSys = %d\n", s.BuckHashSys)
	fmt.Fprintf(fi,"# GCSys = %d\n", s.GCSys)
	fmt.Fprintf(fi,"# OtherSys = %d\n\n", s.OtherSys)

	fmt.Fprintf(fi,"# NextGC = %d\n", s.NextGC)
	fmt.Fprintf(fi, "# LastGC = %d\n", s.LastGC)
	fmt.Fprintf(fi, "# PauseNs = %d\n", s.PauseNs)
	fmt.Fprintf(fi, "# PauseEnd = %d\n", s.PauseEnd)
	fmt.Fprintf(fi, "# NumGC = %d\n", s.NumGC)
	fmt.Fprintf(fi, "# NumForcedGC = %d\n", s.NumForcedGC)
	fmt.Fprintf(fi, "# GCCPUFraction = %v\n", s.GCCPUFraction)
	fmt.Fprintf(fi,"# DebugGC = %v\n", s.DebugGC)

	return nil
}
