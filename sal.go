package main

import (
  "os"
  "io/ioutil"
  "fmt"
)

// id is passed to identify the goRoutine in the println statements
func searchAndLog(file        string,
                  dirpath     chan string,
                  goRoutTermId  chan int,
                  id          int) {
  // For loop ensures that this goRoutine keeps on running as long as
  // main goRoutine is runnign
  for {
    directory, ok := <- dirpath
    if !ok {
      goRoutTermId <- id
      return
    }
    files, _ := ioutil.ReadDir(directory)
    for _, f := range files {
      // Replace ".DS_Store" with something else to search for
      // other file
      if f.Name() == file {
          fmt.Println("goRoutine # ", id, directory, f.Size())
      }
    }
  }
}

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
