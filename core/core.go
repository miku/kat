package core

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const help = `kat - Preview.app for the command line
 _                        _            
(_)                      (_)           
(_)     _   _  _  _    _ (_) _  _      
(_)   _(_) (_)(_)(_) _(_)(_)(_)(_)     
(_) _(_)    _  _  _ (_)  (_)           
(_)(_)_   _(_)(_)(_)(_)  (_)     _     
(_)  (_)_(_)_  _  _ (_)_ (_)_  _(_)    
(_)    (_) (_)(_)(_)  (_)  (_)(_)    

Plain text, directories, PDF, JPG, PNG, gif, MARC, zip, tgz, rar, mp3, odt,
doc, docx, xlsx, tar, tar.gz, dmg, djvu, deb, rpm.

$ kat FILE
`

// Version.
const Version = "0.1.3"

func Run() {
	version := flag.Bool("version", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println(help)
		os.Exit(0)
	}

  if os.Args[1] == "cmd" {
    args := os.Args[2:]
    if len(args) == 0 {
      fmt.Println(help)
      os.Exit(0)
    }

    for _, arg := range args {
      handler, _ := GetCMDHandler(arg)
      b, err := handler(arg)
      if err != nil {
        log.Fatal(err)
      }
      fmt.Println(string(b))
    }
    os.Exit(0)
  }

	for _, arg := range flag.Args() {
		if st, err := os.Stat(arg); err == nil {
			if st.IsDir() {
				out, err := exec.Command("tree", "-sh", arg).Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf(string(out))
			} else {
				handler, _ := GetHandler(arg)
        // fallback to default handler when status is false
				b, err := handler(arg)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf(string(b))
			}
		} else {
      log.Fatal(arg, ": no such file or directory")
    }
	}  
}
