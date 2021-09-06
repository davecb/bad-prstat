package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// prstat is patterned after Solaris prstat, but for Linux needs to have access
//to both the /var/account/pacct file and the output of ps. As the API for pacct
// is non-existent, we use the native commands and parse their output

// pacct is what dump-pacct returns, the stats for all process which have exited today
type pacct struct {
	Command   string
	Version   string
	Utime     float32
	Stime     float32
	Elapsed   float32
	Uid       int
	Gid       int
	Avmem     float32
	Chars     float32
	Pid       int
	Ppid      int
	Flags     string
	Exit      int
	Pty       string
	StartTime time.Time
}

func main() {
	// parse options

	mainLine()
	os.Exit(0)
}

// mainLine does the combination of psacct and ls into a reportable structure
func mainLine() {
	// main loop
	// get a new acct structure into a map
	getPacct()
	// if first, ignore
	// subtract from first giving a new map
	// swap new and old, discard old
	// add ps to new map
	// report
}

// getPacct used dump-acct to get a process accounting file's contents
func getPacct() map[int]bool {
	var sample pacct
	var fred map[int]bool

	out, err := exec.Command("/usr/sbin/dump-acct", "/var/account/pacct").Output()
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", out)
	f := bytes.NewReader(out)
	r := csv.NewReader(f)
	r.Comma = '|'
	r.Comment = '#'
	r.FieldsPerRecord = -1 // ignore differences
	r.LazyQuotes = true    // allow bad quoting
	r.TrimLeadingSpace = true // un-pad beginnings
	for nr := 0; ; nr++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Errorf("error reading csv data, %v", err))
		}
		log.Printf("record = %#v", record)

		sample.Command = strings.TrimSpace(record[0])
		sample.Version = record[1]
		sample.Utime  = toFloat32(record[2])
		sample.Stime = toFloat32(record[3])
		sample.Elapsed = toFloat32(record[4])
		sample.Uid= toInt(record[5])
		sample.Gid = toInt(record[6])
		sample.Avmem = toFloat32(record[7])
		sample.Chars  = toFloat32(record[8])
		sample.Pid  = toInt(record[9])
		sample.Ppid   = toInt(record[10])
		sample.Flags   = strings.TrimSpace(record[11])
		sample.Exit  = toInt(record[12])
		sample.Pty = strings.TrimSpace(record[13])
		sample.StartTime = toTime(record[14])

		log.Printf("sample = %#v\n", sample)
		break
	}
	return fred
}

// toInt returns a normal int from a string that's expected to be correct
// by construction
func toInt(s string) int {
	var i int
	var err error

	if i, err = strconv.Atoi(s); err != nil {
		fmt.Printf("Programmer error: not an int, %T, %v\n", s, s)
		panic(err)
	}
	return i
}

// toFloat32 does the same with floats
func toFloat32(s string) float32 {
	var f float64
	var err error

	if f, err = strconv.ParseFloat(s, 32); err != nil {
		fmt.Printf("Programmer error: not an float, %T, %v\n", s, s)
		panic(err)
	}
	return float32(f)
}

time.Parse(layout value) time err
// toTime does the same with time strings
func toTime(s string) time.Time {
	var t time.Time
	var err error

	if f, err = strconv.ParseFloat(s, 32); err != nil {
		fmt.Printf("Programmer error: not an float, %T, %v\n", s, s)
		panic(err)
	}
	return float32(f)
}