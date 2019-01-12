package minotaur

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

// Run watches for events on the filesystem and runs a command when they happen.
func Run(args []string, _ io.Reader) error {
	args = args[1:]

	c, err := parseFlags(args)
	if err != nil {
		return errors.Wrap(err, "parseFlags")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "fsnotify.NewWatcher")
	}
	defer watcher.Close()

	done := make(chan bool)
	var timeout <-chan time.Time
	events := make(map[string]bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					stat, err := os.Stat(event.Name)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Couldn't stat created thing: %s\n", err)
					} else if stat.IsDir() && c.include.MatchString(event.Name) && !c.ignore.MatchString(event.Name) {
						err := watcher.Add(event.Name)
						if err != nil {
							fmt.Fprintf(os.Stderr, "failed to watch %s: %s\n", event.Name, err)
						} else if c.verbose {
							fmt.Fprintf(os.Stderr, "watching %s\n", event.Name)
						}
					}
				}

				events[event.Op.String()+"\t"+event.Name] = true
				timeout = time.After(time.Second)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			case <-timeout:
				s := make([]string, 0, len(c.script)+len(events))
				s = append(s, c.script...)
				for e := range events {
					s = append(s, e)
				}
				events = make(map[string]bool)
				cmd := exec.Command(s[0], s[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil && c.verbose {
					fmt.Fprintf(os.Stderr, "script (%q) failed: %s\n", s, err)
				}
			}

		}
	}()

	for _, path := range c.dirs {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				return nil
			}

			if c.ignore.MatchString(path) {
				return filepath.SkipDir
			}
			if !c.include.MatchString(path) {
				return nil
			}

			if c.verbose {
				fmt.Fprintln(os.Stderr, "watching "+path)
			}
			return errors.Wrap(watcher.Add(path), "fsnotify.Watcher.Add")
		})
		if err != nil {
			return errors.Wrap(err, "filepath.Walk")
		}
	}
	<-done

	return nil
}
