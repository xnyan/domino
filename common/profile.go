package common

import (
	"os"
	"runtime"
	"runtime/pprof"
)

func example() {
	destDir := "./fp-exp/"
	Profile(destDir)
	HeapProfile(destDir + "mem.log")
}

func Profile(destDir string) {
	nList := []string{"heap", "goroutine", "allocs", "threadcreate", "block", "mutex"}
	for _, l := range nList {
		f := destDir + "/" + l + ".log"
		DumpProfile(l, f, 0)
	}
}

func HeapProfile(f string) {
	mf, err := os.Create(f)
	if err != nil {
		logger.Fatal("could not create memory profile: ", err)
	}
	defer mf.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(mf); err != nil {
		logger.Fatal("could not write memory profile: ", err)
	}
}

func DumpProfile(label, f string, debug int) {
	mf, err := os.Create(f)
	if err != nil {
		logger.Fatalf("could not create %s profile: %v", label, err)
	}
	defer mf.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.Lookup(label).WriteTo(mf, debug); err != nil {
		logger.Fatal("could not write %s profile: %v", label, err)
	}
}
