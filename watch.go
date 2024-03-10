package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func createWatcher(ctx context.Context, pathToWatch string, watcherNotifications chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
					watcherNotifications <- struct{}{}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-ctx.Done():
				watcher.Close()
				return
			}
		}
	}()

	// Add a path.
	directories := findAllDirectories(pathToWatch)
	fmt.Printf("Watching directories: %+v\n", directories)

	for _, dr := range directories {
		err = watcher.Add(dr)
		if err != nil {
			return err
		}
	}

	return nil
}

func findAllDirectories(pathToWatch string) []string {
	directories := []string{}
	filepath.WalkDir(pathToWatch, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			directories = append(directories, path)
		}

		return nil
	})

	return directories
}
