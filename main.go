package main

// A simple multithreaded program to calculate how many
// and total disk space occupyed by all ".DS_Store" files.
// It also logs the location of each file along with its
// size

import (
  "fmt"
  "os"
  "flag"
  "path/filepath"
)

func main() {

  // Command line flag to pass in how many goRouts to swapn
  var numGoRout  int
  var directory   string
  var file        string

  flag.StringVar(&file, "f", "", "file to search for")
  flag.IntVar(&numGoRout, "t", 3, "number of goRoutines")
  flag.StringVar(&directory, "d", "", "directory to scan")

  flag.Parse()

  if file == "" {
    fmt.Println("-f flag is required")
    os.Exit(2)
  }

  if directory == "" {
    fmt.Println("-d flag is required")
    os.Exit(2)
  }

  // 5 Channels for 5 goRoutines
  var dir [5]chan string
  for i := 0; i < numGoRout; i++ {
    dir[i] = make(chan string)
  }

  // Using this channel goRoutines will communicate back to main thread
  // when they are about to exit
  var goRoutTermId = make(chan int)

  // This is the function that will be passed to filepath.Walk()
  // "select" will be executed only if path points to directory
  walkFunc := func(path string, info os.FileInfo, err error) error {
    if isDirectory(path) {
      select {
        case dir[0] <- path:
        case dir[1] <- path:
        case dir[2] <- path:
        case dir[3] <- path:
        case dir[4] <- path:
        }
    }
    return nil
  }

  // Create numGoRout goRoutines of searchAndLog()
  for i := 0; i < numGoRout; i++ {
    go searchAndLog(file, dir[i], goRoutTermId, i+1)

  }

  filepath.Walk(directory, walkFunc)

  // Close channels so that goRoutines can terminate themselves
  for i := 0; i < numGoRout; i++ {
    close(dir[i])
  }

  // Wait for all goRoutines to be terminated before ending main
  for i := 0; i < numGoRout; i++ {
    _ = <- goRoutTermId
  }
}
