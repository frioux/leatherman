package minotaur

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func Run(args []string, _ io.Reader) error {
	args = args[1:]

	c, err := parseFlags(args)
	if err != nil {
		return fmt.Errorf("parseFlags: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fsnotify.NewWatcher: %w", err)
	}
	defer watcher.Close()

	for _, path := range c.dirs {
		if err := addDir(watcher, c, path); err != nil {
			return err
		}
	}

	var timeout <-chan time.Time
	events := make(map[string]bool)

	if !c.noRunAtStart {
		timeout = time.After(0)
	}

LOOP:
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("watcher went away")
			}

			// sink the ship if a root disappears
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				for _, path := range c.dirs {
					if path == event.Name {
						return errors.New("deleted root, capsizing")
					}
				}
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				stat, err := os.Stat(event.Name)
				if err != nil {
					if os.IsNotExist(err) {
						continue LOOP
					}
					fmt.Fprintf(os.Stderr, "Couldn't stat created thing: %s\n", err)
				} else if stat.IsDir() {
					err := addDir(watcher, c, event.Name)
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
				return errors.New("watcher went away")
			}
			fmt.Println("error:", err)
		case <-timeout:
			s := make([]string, 0, len(c.script)+len(events))
			s = append(s, c.script...)
			if c.includeArgs {
				for e := range events {
					s = append(s, e)
				}
			}
			events = make(map[string]bool)
			if c.report {
				fmt.Println("==============", time.Now().Format("2006-01-02 03:04:05"), "==============")
			}

			cmd := exec.Command(s[0], s[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil && c.verbose {
				fmt.Fprintf(os.Stderr, "script (%q) failed: %s\n", s, err)
			}
			if c.report {
				fmt.Println("=================================================")
			}
		}

	}
}

func addDir(watcher *fsnotify.Watcher, c config, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
		if err := watcher.Add(path); err != nil {
			return fmt.Errorf("fsnotify.Watcher.Add: %w", err)
		}
		return nil
	})
}
