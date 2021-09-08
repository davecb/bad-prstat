run: /home/davecb/go/bin/prstat
	sudo /home/davecb/go/bin/prstat

/home/davecb/go/bin/prstat: *.go
	go install 
