package main

// A simple multithreaded program to calculate how many
// and total disk space occupyed by all ".DS_Store" files.
// It also logs the location of each file along with its
// size

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {

	// Command line flag to pass in how many goRouts to swapn
	var numGoRout int
	var directory string
	var file string

	flag.StringVar(&file, "f", "", "file to search for")
	flag.IntVar(&numGoRout, "t", 3, "number of goRoutines")
	flag.StringVar(&directory, "d", "", "directory to scan")

	flag.Parse()

	runtime.GOMAXPROCS(numGoRout)

	if file == "" {
		fmt.Println("-f flag is required")
		os.Exit(2)
	}

	if directory == "" {
		fmt.Println("-d flag is required")
		os.Exit(2)
	}

	e, err := os.OpenFile("dsLog.log", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening dsLog.log file: ", err)
		os.Exit(1)
	}
	defer e.Close()

	dsLogger := log.New(e, "", log.Ldate|log.Ltime)

	// Set buffer length to 5 so that directory walker function does not to
	// wait for goRoutines to accept directory
	dirs := make(chan string, 5)

	infoLog := make(chan dsDetails, 10)

	var wg sync.WaitGroup

	// This is the function that will be passed to filepath.Walk()
	// "select" will be executed only if path points to directory
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if isDirectory(path) {
			dirs <- path
		}
		return nil
	}

	go func(il <-chan dsDetails) {
		defer wg.Done()
		var totalSize int64 = 0
        var totalFiles int = 0
		for {
			d, ok := <-il
			if !ok {
				dsLogger.Println("Total Size: ",
					float64(totalSize)/1024.0, "Kb")
                dsLogger.Println("Total Number of Files: ", totalFiles)
				return
			}
            totalFiles += 1
			totalSize += d.size
			dsLogger.Println(d.info)
		}
	}(infoLog)
	wg.Add(1)

	// Create numGoRout goRoutines of searchAndLog()
	for i := 0; i < numGoRout; i++ {
		wg.Add(1)
		go searchAndLog(file, dirs, infoLog, &wg, i+1)
	}

	filepath.Walk(directory, walkFunc)

	// Close channels so that goRoutines can terminate themselves
	close(dirs)
	close(infoLog)

	// Wait for all goRoutines to be terminated before ending main
	wg.Wait()
}
