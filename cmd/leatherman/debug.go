package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"runtime/trace"

	_ "net/http/pprof"
)

func startDebug() {
	port := os.Getenv("LMHTTPPROF")
	if port == "" {
		port = "6060"
	}
	go func() {
		http.ListenAndServe("localhost:"+port, nil)
	}()
	if os.Getenv("LMTRACE") != "" {
		fh, err := os.Create(os.Getenv("LMTRACE"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open LMTRACE (%s): %s\n", os.Getenv("LMTRACE"), err)
			os.Exit(1)
		}
		trace.Start(fh)
	}

	if os.Getenv("LMPROF") != "" {
		fh, err := os.Create(os.Getenv("LMPROF"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open LMPROF (%s): %s\n", os.Getenv("LMPROF"), err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(fh)
	}
}

func stopDebug() {
	if os.Getenv("LMTRACE") != "" {
		trace.Stop()
	}

	if os.Getenv("LMPROF") != "" {
		pprof.StopCPUProfile()
	}
}
