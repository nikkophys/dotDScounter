package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type dsDetails struct {
	info string
	size int64
}

// id is passed to identify the goRoutine in the println statements
func searchAndLog(
	file string,
	dirpath <-chan string,
	logData chan<- dsDetails,
	wg *sync.WaitGroup,
	id int) {

	defer wg.Done()
	ld := dsDetails{}

	// For loop ensures that this goRoutine keeps on running as long as
	// main goRoutine is runnign
	for {
		directory, ok := <-dirpath
		if !ok {
			return
		}
		files, _ := ioutil.ReadDir(directory)
		for _, f := range files {
			// Replace ".DS_Store" with something else to search for
			// other file
			if f.Name() == file {
				ld.size = f.Size()
				ld.info = fmt.Sprintf("goRoutine # %d, %s, %d",
					id, directory, f.Size())
				logData <- ld
				// fmt.Println("goRoutine # ", id, directory, f.Size())
			}
		}
	}
}

// Returns true if path is a directory otherwise false
func isDirectory(path string) bool {

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("ERROR >>> ", err)
		return false
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println("ERROR >>> ", err)
		return false
	}

	return (fi.Mode()).IsDir()
}
