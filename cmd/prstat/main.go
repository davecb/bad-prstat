package main

import "os"

// prstat is patterned after Solaris prstat, but for Linux needs to have access
//to both the /var/account/pacct file and the output of ps. As the API for pacct
// is non-existent, we use the native commands and parse their output

import (
	"fmt"
)


func main() {
	fmt.Printf("acctcom\n")
	fp, err := os.Open("/var/account/pacct")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	os.Exit(0)
}
