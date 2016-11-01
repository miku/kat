package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		if st, err := os.Stat(arg); err == nil {
			if st.IsDir() {
				out, err := exec.Command("ls", arg).Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf(string(out))
			} else {
				out, err := exec.Command("cat", arg).Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf(string(out))
			}
		}
	}
}
