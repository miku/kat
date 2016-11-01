SHELL = /bin/bash

kat: cmd/kat/main.go
	go build -o kat cmd/kat/main.go

install: kat
	cp kat $(HOME)/bin

clean:
	rm -f kat

