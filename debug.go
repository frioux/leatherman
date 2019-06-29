package main

import (
	"fmt"
	"net"
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
			if oerr, ok := err.(*net.OpError); ok && oerr.Op == "listen" {
				if ierr, ok := oerr.Err.(*os.SyscallError); ok && ierr.Err.Error() == "address already in use" {
					err := http.ListenAndServe("localhost:0", nil)
					if err != nil {
						fmt.Fprintf(os.Stderr, "failed to http.ListenAndServe: %s\n", err)
					}
				}
			} else {
				fmt.Fprintf(os.Stderr, "failed to http.ListenAndServe: %s\n", err)
			}
		}
	}()
	if path := os.Getenv("LMTRACE"); path != "" {
		fh, err := os.Create(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open LMTRACE (%s): %s\n", path, err)
			os.Exit(1)
		}
		err = trace.Start(fh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to trace.Start: %s\n", err)
		}
	}

	if path := os.Getenv("LMPROF"); path != "" {
		fh, err := os.Create(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open LMPROF (%s): %s\n", path, err)
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
