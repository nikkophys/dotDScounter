package main

// A simple multithreaded program to calculate how many
// and total disk space occupyed by all ".DS_Store" files.
// It also logs the location of each file along with its
// size

import (
  "fmt"
  "io/ioutil"
  "os"
  "flag"
  "path/filepath"
)

// Returns true if path is a directory otherwise false
func isDirectory(path string) bool {

  file, err := os.Open(path)
  if err != nil {
    fmt.Println(err)
    return false
  }
  defer file.Close()

  fi, err := file.Stat()
  if err != nil {
    fmt.Println(err)
    return false
  }

  return (fi.Mode()).IsDir()
}

// id is passed to identify the thread in the println statements
func searchAndLog(file string, dirpath chan string, id int) {
  // For loop ensures that thread keeps on running as long as
  // main thread is runnign
  for {
    directory := <- dirpath
    files, _ := ioutil.ReadDir(directory)
    for _, f := range files {
      // Replace ".DS_Store" with something else to search for
      // other file
      if f.Name() == file {
        // fmt.Println("Thread # ", id + 1, directory, f.Size())
      }
    }
  }
}

func main() {

  // Command line flag to pass in how many threads to swapn
  var numThreads  int
  var directory   string
  var file        string

  flag.StringVar(&file, "f", "", "file to search for")
  flag.IntVar(&numThreads, "t", 3, "number of threads")
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

  // 5 Channels for 5 threads
  var dir [5]chan string
  for i := 0; i < numThreads; i++ {
    dir[i] = make(chan string)
  }

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

  // Create numThreads threads of searchAndLog()
  for i := 0; i < numThreads; i++ {
    go searchAndLog(file, dir[i], i)

  }

  go filepath.Walk(directory, walkFunc)

  var input string
  fmt.Scanln(&input)
}
