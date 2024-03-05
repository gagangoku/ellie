package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	trace2 "runtime/trace"

	"github.com/felixge/fgprof"
	"github.com/samber/lo"
)

func _pretty(table Table, sep string, quote bool) []string {
	list := make([]string, 0, table.nRows+1)
	list = append(list, strings.Join(table.colNames, sep))

	for i := 0; i < table.nRows; i++ {
		row := make([]string, 0, len(table.colNames))
		for _, col := range table.cols {
			if quote {
				row = append(row, "\""+interfaceToString(col[i])+"\"")
			} else {
				row = append(row, interfaceToString(col[i]))
			}
		}
		list = append(list, strings.Join(row, sep))
	}
	return list
}

func trimLines(str string) string {
	lines := strings.Split(str, "\n")
	lines = lo.Map(lines, func(item string, index int) string { return strings.TrimSpace(item) })
	lines = lo.Filter(lines, func(item string, index int) bool { return item != "" })
	return strings.Join(lines, "\n")
}

func setupProfiling() CleanupFn {
	// CPU profiling
	profileDirFlag := "/tmp/ell-pprof"

	if _, err := os.Stat(profileDirFlag); os.IsNotExist(err) {
		os.Mkdir(profileDirFlag, os.ModePerm)
	}

	enableRuntimeTraceFlag := true
	cpuProfF, _ := os.Create(fmt.Sprintf("%s/cpu.prof", profileDirFlag))
	pprof.StartCPUProfile(cpuProfF)

	memProfF, _ := os.Create(fmt.Sprintf("%s/mem.out", profileDirFlag))

	// Set block profiling rate
	runtime.SetBlockProfileRate(100_000_000) // WARNING: Can cause some CPU overhead

	var traceFile *os.File = nil

	// go tool trace
	if enableRuntimeTraceFlag {
		f, err := os.Create(fmt.Sprintf("%s/trace.out", profileDirFlag))
		if err != nil {
			panic(err)
		}
		traceFile = f
		err = trace2.Start(f)
		if err != nil {
			panic(err)
		}
	}

	// fgprof
	http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
	go func() {
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()

	return func() {
		pprof.StopCPUProfile()
		cpuProfF.Close()

		pprof.WriteHeapProfile(memProfF)
		memProfF.Close()

		if enableRuntimeTraceFlag {
			trace2.Stop()
			traceFile.Close()
		}

		// Save block profile
		if profileDirFlag != "" {
			file, _ := os.Create(fmt.Sprintf("%s/block.prof", profileDirFlag))
			err := pprof.Lookup("block").WriteTo(file, 0)
			if err != nil {
				fmt.Println("ERROR: failed to write block profile:", err)
			}
		}
	}
}
