package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	pathToMonitor := "."
	mainFile := "main.go"

	if len(os.Args) > 1 {
		pathToMonitor = os.Args[1]
		if len(os.Args) > 2 {
			mainFile = os.Args[2]
		}
	}

	ctx := context.Background()
	watcherNotifications := make(chan struct{})
	err := createWatcher(ctx, pathToMonitor, watcherNotifications)
	if err != nil {
		fmt.Printf("Error occured: %s\n", err)
		os.Exit(-1)
	}

	cmd, err := startProcess(mainFile)
	if err != nil {
		fmt.Printf("Error occured when running program: %s\n", err)
		os.Exit(-1)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-sigc
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		os.Exit(0)
	}()

	// Wait for watcher events
	for range watcherNotifications {
		//Stop
		log.Println("Stopping process")
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

		log.Println("Restaring process")
		cmd, err = startProcess(pathToMonitor)
		if err != nil {
			fmt.Printf("Unable to restart the process %s\n", err)
			os.Exit(-1)
		}

		log.Println("Restarted process")
	}
	// Stop existing process and start new process
}
