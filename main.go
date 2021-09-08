package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// prstat is patterned after Solaris prstat, but for Linux needs to have access
//to both the /var/account/ file and the output of ps. As the APIs for each are
// non-existent, we use the native commands and parse their output

// Pacct what dump-acct returns, the stats for all process which have exited today
type Pacct struct {
	Command string // short name of command
	//Version   string // uninteresting, unless it's not v3
	Utime   float32 // user time in ticks (s/100) or seconds, 0 if unknown
	Stime   float32 // system time in ticks (s/100) or seconds, 0 if unknown
	CpuTime float32 // sum of the above
	Elapsed float32 // elapsed time in seconds
	Uid     int
	Gid     int
	Avmem   float32 // average memory usage (kB)
	//Chars     float32 // unused
	Pid  int
	Ppid int
	//Flags     string  // uninteresting
	//Exit      int     // uninteresting
	//Pty       string  // uninteresting
	StartTime time.Time
}


type options struct {
	Verbose   bool
	Debug     bool
	SortOrder string
	ProcList string
}
var debug bool

// main parses options and starts mainline
func main() {
	var opts options
	var interval, count int

	flag.StringVar(&opts.ProcList,"p", "", "specify a comma-separated list of processes")
	flag.StringVar(&opts.SortOrder, "s", "cpu", "specify sort order, either cpu or mem")
	flag.BoolVar(&opts.Verbose, "v", false, "turn verbose reporting on")
	flag.BoolVar(&opts.Debug, "d", false, "turn debug reporting on")
	flag.Parse()

	switch len(flag.Args()) {
	case 0:
		// defaults
		interval = 60
		count = 1
	case 1:
		interval = parseArg(flag.Arg(0), "interval")
		count = 1
	case 2:
		interval = parseArg(flag.Arg(0), "interval")
		count = parseArg(flag.Arg(1), "count")
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized arguement in %q\n", flag.Args())
		usage()
		os.Exit(1)
	}
	debug = opts.Debug

	mainLine(interval, count, opts)
	os.Exit(0)
}



func usage() {
	fmt.Fprintf(os.Stderr, "Usage: prstat [-v, -d, -s cpu|mem, -p proclist] interval count\n")
	fmt.Fprintf(os.Stderr, "       default inteval is 60 seconds, count is 1\n")
}

func parseArg(s, name string) int {
	var i int
	var err error

	if i, err = strconv.Atoi(s); err != nil {
		fmt.Fprintf(os.Stderr,	"Programmer error: %s was not an int, halting. %s=%q\n", name, name, s)
		panic(err)
		//os.Exit(1)
	}
	return i
}

// Dprintf is really debug.Printf
func Dprintf(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v)
	}
}