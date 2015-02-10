package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Name) && !ev.IsAttrib() {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			}
		}
	}()

	watcherLog("Watching %s", path)
	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func watch() {
	watchPath, err := filepath.Abs(watchPath())
	if err != nil {
		fatal(err)
	}

	filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if (len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")) || isExcludedDir(path) {
				return filepath.SkipDir
			}

			watchFolder(path)
		}

		return err
	})
}
