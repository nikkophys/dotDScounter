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

func main() {

  // Command line flag to pass in how many threads to swapn
  var numThreads int
  flag.IntVar(&numThreads, "t", 5, "number of threads")
  flag.Parse()
  fmt.Println("numThreads: ", numThreads)

  // 5 Channels for 5 threads
  var dir [5]chan string
  for i, _ := range dir {
    dir[i] = make(chan string)
  }

  // This is the function that will be passed to filepath.Walk()
  // "select" will be executed only if path points to directory
  walkFunc := func(path string, info os.FileInfo, err error) error {
    fmt.Println("Visited: ", path)
    if isDirectory(path) {
      select {
        case dir[0] <- path:
          fmt.Println("Thread: 1")
        case dir[1] <- path:
          fmt.Println("Thread: 2")
        case dir[2] <- path:
          fmt.Println("Thread: 3")
        case dir[3] <- path:
          fmt.Println("Thread: 4")
        case dir[4] <- path:
          fmt.Println("Thread: 5")
        }
    }
    return nil
  }

  // Create 5 threads of searchAndLog()
  for i := 0; i < numThreads; i++ {
    go func(i int) {
      for {
        searchAndLog(dir[i], i)
      }
    }(i)
  }

  go filepath.Walk("/Users/nikhil/Workspace/Test", walkFunc)

  var input string
  fmt.Scanln(&input)
}

// id is passed to identify the thread in the println statements
func searchAndLog(dirpath chan string, id int) {
  directory := <- dirpath
  files, _ := ioutil.ReadDir(directory)
  for _, f := range files {
    if f.Name() == ".DS_Store" {
      fmt.Println("########################")
      fmt.Println("Thread # ", id + 1, directory, f.Size())
    }
  }
}
