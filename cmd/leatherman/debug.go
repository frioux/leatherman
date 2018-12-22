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
		err := http.ListenAndServe("localhost:"+port, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to http.ListenAndServe: %s\n", err)
		}
	}()
	if os.Getenv("LMTRACE") != "" {
		fh, err := os.Create(os.Getenv("LMTRACE"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open LMTRACE (%s): %s\n", os.Getenv("LMTRACE"), err)
			os.Exit(1)
		}
		err = trace.Start(fh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to trace.Start: %s\n", err)
		}
	}

	if os.Getenv("LMPROF") != "" {
		fh, err := os.Create(os.Getenv("LMPROF"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open LMPROF (%s): %s\n", os.Getenv("LMPROF"), err)
			os.Exit(1)
		}
		err = pprof.StartCPUProfile(fh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to pprof.StartCPUProfile: %s\n", err)
		}
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
