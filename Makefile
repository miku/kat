SHELL = /bin/bash

kat: cmd/kat/main.go
	go build -o github.com/miku/kat

install: kat
	cp kat $(HOME)/bin

clean:
	rm -f kat

