## dotDScounter

A multi threaded command line utility, that shows you how much of disk space is occupyed by ".DS_Store" files under a given directory tree.

## Instructions

Using a command line flag __-t=4__ you can specify number of threads that this utility will spawn. This utility allows up to 6 threads to be spawned. This is an arbitary restriction.

`$ ./dotDScounter -d=/`

Above example spawn 6 threads. One of which will do directory walk, passing each directory it encounters to one of the four listener threads which will search for **.DS_Store** files in those directories. 6th thread is the logger thread. It recieves messages form the 4 searching threads and logs those messages to files as well as print them to screen.

## Tests


##License

``dotDScounter`` is released as open-source software under the GNU General Public License (GPL), version 2 or later.

## About
