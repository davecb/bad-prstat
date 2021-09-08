package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// getPs uses the ps command  to get a snapshot of processes which haven't exited yet
func getPs(opts options) map[int]Pacct {
	var psMap = make(map[int]Pacct)

	if opts.Verbose {
		fmt.Printf("#Command, CpuTime, Elapsed, Uid, Gid, Avmem, Pid, Ppid, StartTime\n")
	}


	out, err := exec.Command("/usr/bin/ps", "-e", "-o", "comm,cputimes,etimes,uid,gid,vsize,pid,ppid,stime").Output()
	if err != nil {
		panic(err)
	}
	Dprintf("out = \n%s\n", out)
	f := bytes.NewReader(out)
	r := csv.NewReader(f)
	r.Comma = ' '
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
		Dprintf("record = %#v", record)

		if sample := parsePsRecord(record); sample.Command != "COMMAND" {
			// ignore headers, just store data lines
			psMap[sample.Pid] = sample

			if opts.Verbose {
				fmt.Printf("%s, %.1f, %.1f, %d, %d, %.0f, %d, %d, %s\n",
					sample.Command, sample.CpuTime, sample.Elapsed, sample.Uid, sample.Gid,
					sample.Avmem, sample.Pid, sample.Ppid, sample.StartTime.Format(time.RFC3339))
			}
		}
	}
	return psMap

}

// parsePsRecord parses a line from ps.
func parsePsRecord(record []string) Pacct {
	var sample Pacct

	if sample.Command = strings.TrimSpace(record[0]); sample.Command == "COMMAND" {
		return sample // It's just a header line
	}

	sample.CpuTime = toFloat32(record[1]) // seconds
	sample.Stime = 0
	sample.Utime = 0
	sample.Elapsed = toFloat32(record[2]) // seconds
	sample.Uid = toInt(record[3])
	sample.Gid = toInt(record[4])
	sample.Avmem = toFloat32(record[5]) // vsize: address space in kb
										// similar to "avmem" in psacct
	sample.Pid = toInt(record[6])
	sample.Ppid = toInt(record[7])
	sample.StartTime = hourToTime(record[8])

	Dprintf("sample = %#v\n", sample)
	Dprintf("sample2 = %v\n", sample)
	return sample
}

// hourToTime will work ONLY IF the processing is done the same day.
// True by construction, but stupid.
func hourToTime(s string) time.Time {
	var now = time.Now()
	var hhmm = strings.Split(s, ":")

	hh := toInt(hhmm[0])
	mm := toInt(hhmm[1])
	y,m, d := now.Date()
	then := time.Date(y, m, d, hh, mm, 0, 0, time.Local)

	//Dprintf("s = %s, now = %v,  then = %v\n", s, now, then)
	return then
}