package main

// old.go -- this was an early attempt to parse psacct files in native go
// For prstat, we use dump-acct and ls, which know all the tricks (;-))
import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"
)

const (
	ACCT_COMM = 16

	// Bits that may be set in ac_flag field
	AFORK = 0x01 // Has executed fork, but no exec
	ASU   = 0x02 // Used superuser privileges
	ACORE = 0x08 // Dumped core
	AXSIG = 0x10 // Killed by a signal
)

type compT uint16
//The comp_t data type is a floating-point value consisting of a
// 3-bit, base-8 exponent, and a 13-bit mantissa.  A value, c, of
// this type can be converted to a (long) integer as follows:
// v = (c & 0x1fff) << (((c >> 13) & 0x7) * 3);
// The Acutime, Acstime, and Acetime fields measure time in
// "clock ticks"; divide these values by sysconf(_SC_CLK_TCK) to
// convert them to seconds.

type acct struct {
	Flag     byte            // Accounting flags
	Version  byte            // Always set to ACCT_VERSION (3)
	Tty      uint16          // Controlling terminal
	ExitCode uint32          // Process termination status, see wait(2)
	Uid      uint32          // Accounting user ID
	Gid      uint32          // Accounting group ID
	Pid      uint32          // Process ID
	Ppid     uint32          // Parent process ID
	Btime    uint32          // Process creation time, ticks
	Etime    float32         // Elapsed time, s
	Utime    compT           // User CPU time, ticks
	Stime    compT           // System time, ticks
	Mem      compT           // Average memory usage (kB)
	Io       uint16          // Characters transferred (unused)
	Rw       uint16          // Blocks read or written (unused)
	Minflt   uint16          // Minor page faults
	Majflt   uint16          // Major page faults
	Swaps    uint16          // Number of swaps (unused)
	Comm     [ACCT_COMM]byte // Command name
}

func oldMain() {
	fmt.Printf("acctcom\n")
	fp, err := os.Open("/var/account/Pacct")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	fmt.Printf("#cmd, uid, gid, pid, ppid, AFORK, ASU, ACORE, AXSIG, version, tty, exit, " +
		"btime, etime, utime, stime, mem, io, rw, minf, majf, swaps\n")

	var a acct
	log.Printf("size = %d", unsafe.Sizeof(a))
	for {


		err := binary.Read(fp, binary.LittleEndian, &a)
		if err == io.EOF{
			break
		}

		cmd := strings.ReplaceAll(string(a.Comm[:]), "\x00", "")
		//fmt.Printf("accounting record for %q = %#v\n", cmd, a)
		if a.Version != 3 {
			// FIXME test for a programmer error: mine
			panic("version is not a 3")
		}
		// Print each record as csv
		fmt.Printf("%s, ", cmd)
		fmt.Printf("%d, %d, %d, %d, ", a.Uid, a.Gid, a.Pid, a.Ppid)
		fmt.Printf("%s, ", flagToString(a.Flag))
		fmt.Printf("%d, %d, %d, ", a.Version, a.Tty, a.ExitCode)
		fmt.Printf("%#x, ", a.Btime) // Time since boot
		fmt.Printf("%fs, ", a.Etime/100) // elapsed time
		fmt.Printf("%s, ", time.Duration(compTtoDuration(a.Utime)))
		fmt.Printf("%s, ", time.Duration(compTtoDuration(a.Stime)))
		fmt.Printf("%d, 0, 0, %d, %d, 0, ", a.Mem, a.Minflt, a.Majflt)

		fmt.Printf("\n")
	}

}

// flagToString formats the flags byte as four csv values
func flagToString(f byte) string {
	var flag string

	if f & AFORK == AFORK {
		flag = "AFORK, "
	} else {
		flag = ", "
	}
	if f & ASU == ASU {
		flag += "ASU, "
	} else {
		flag += ", "
	}
	if f & ACORE == ACORE {
		flag += "ACORE, "
	} else {
		flag += ", "
	}
	if f & AXSIG == AXSIG {
		flag += "AXSIG"
	}
	return flag
}

// compTtoDuration supposedly converts a compT to a int64 duration
//
// The comp_t data type is a floating-point value consisting of a
// 3-bit, base-8 exponent, and a 13-bit mantissa.  A value, c, of
// this type can be converted to a (long) integer as follows:
// v = (c & 0x1fff) << (((c >> 13) & 0x7) * 3);
//
// The Acutime, Acstime, and Acetime fields measure time in
// "clock ticks"; divide these values by sysconf(_SC_CLK_TCK) to
// convert them to seconds.
func compTtoDuration(c compT) float64 {
	//var v compT
	var mantissa int16
	var f float64

	//log.Printf("c = %d\n" ,c)
	mantissa = int16(c)
	mantissa = mantissa & 0x1fff // remove exponent
	//log.Printf("mantissa = %d\n", mantissa)
	shift := ((c >> 13) & 0x7) * 3 // turn exponent into a shift amount
	//log.Printf("shift = %d\n", shift)
	f = float64(int64(mantissa) << shift)
	//log.Printf("f = (mantissa << shift) = %f\n", f)
	f /= 100.0 // clock_ticks per second
	//log.Printf(" f /= 100 = %f\n", f)

	return f
	//v = (c & 0x1fff) << (((c >> 13) & 0x7) * 3)
	//v /= 100
	//return int64(v)
}
