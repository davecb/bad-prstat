package main


// mainLine does the coordination
func mainLine(interval, count int, opts options) {
	Dprintf("mainline: interval = %d, count = %d\n", interval, count)
	Dprintf("          opts = %#v\n", opts)

	// main loop
	// get a new acct structure into a map
	//		getPacct(opts)
	// add ps to new map
	      getPs(opts)
	// if first, ignore
	// subtract from first giving a new map
	// swap new and old, discard old

	// report
}

