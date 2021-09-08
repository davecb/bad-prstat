package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//// pacct file
// getPacct uses dump-acct to get a process accounting file's contents
func getPacct() map[int]Pacct {
	var pacctMap = make(map[int]Pacct)

	out, err := exec.Command("/usr/sbin/dump-acct", "/var/account/").Output()
	if err != nil {
		panic(err)
	}
	Dprintf("%s\n", out)
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
		Dprintf("record = %#v", record)

		sample := parsePacctRecord(record)
		pacctMap[sample.Pid] = sample
	}
	return pacctMap
}

// parsePacctRecord parses a process accounting record
func parsePacctRecord(record []string) Pacct {
	var sample Pacct

	sample.Command = strings.TrimSpace(record[0])
	version := record[1]
	if version != "v3" {
		log.Fatalf("The version of record is not v3, halting. version = %s\n", version)
	}
	sample.Utime = toFloat32(record[2])
	sample.Stime = toFloat32(record[3])
	sample.CpuTime = sample.Stime + sample.Stime
	sample.Elapsed = toFloat32(record[4])
	sample.Uid = toInt(record[5])
	sample.Gid = toInt(record[6])
	sample.Avmem = toFloat32(record[7])
	sample.Pid = toInt(record[9])
	sample.Ppid = toInt(record[10])
	sample.StartTime = toTime(record[14])

	Dprintf("sample = %#v\n", sample)
	Dprintf("sample2 = %v\n", sample)
	return sample
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

// toTime does the same with time strings
// Input format is Mon Sep  6 09:40:43 2021
func toTime(s string) time.Time {
	var t time.Time
	var err error

	if t, err = time.ParseInLocation(time.ANSIC, s, time.Local); err != nil {
		fmt.Printf("Programmer error: not an time, %T, %v\n", s, s)
		panic(err)
	}
	return t
}

