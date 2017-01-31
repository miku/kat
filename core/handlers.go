package core

import (
  "os/exec"
  "fmt"
)

const ARGIDENT = "{|ARG|}"

// TODO
// 1.improve this, probably a graph is suited as serveral extentions
// are refering to the same command
// 2.add support for remote server ( send os and file extention,
//  fetch appropirate recipe and register it in runtime and save it )
var defaultHandlers map[string][]string = map[string][]string {
  ".pdf"    : []string{"pdftotext", "-enc", "UTF-8", ARGIDENT, "-" },
  ".jpg"    : []string{"catimg", "-w", "192"                       },
  ".png"    : []string{"catimg", "-w", "192"                       },
  ".gif"    : []string{"catimg", "-w", "192"                       },
  ".mrc"    : []string{"catimg", "-w", "192"                       },
  ".zip"    : []string{"unzip", "-l"                               },
  ".tgz"    : []string{"tar", "tf"                                 },
  ".tar.gz" : []string{"tar", "tf"                                 },
  ".mp3"    : []string{"mp3info", "-x"                             },
  ".rar"    : []string{"unrar", "l"                                },
  ".odt"    : []string{"docd", "-input"                            },
  ".docx"   : []string{"docd", "-input"                            },
  ".doc"    : []string{"antiword"                                  },
  ".xlsx"   : []string{"xlsx2tsv.py"                               },
  ".dmg"    : []string{"hdiutil", "imageinfo"                      },
  ".djvu"   : []string{"djvutxt"                                   },
  ".deb"    : []string{"dpkg", "-c"                                },
  ".rpm"    : []string{"rpm", "-qplv"                              },
}

// View is the generic via cat.
func DefaultHandler(f string) ([]byte, error) {
	if _, err := exec.LookPath("cat"); err != nil {
		return nil, fmt.Errorf("cat is required")
	}
	return exec.Command("cat", f).Output()
}

// Create and return a new handler
func Create(command string, args ...string) Handler{
  newHandler := func(f string) ([]byte, error) {
    if _, err := exec.LookPath(command); err != nil {
      return nil, fmt.Errorf("%s is required", command)      
    }
    var argv []string = make([]string, len(args))
    reparg := false
    for index, arg := range args {
      if arg == ARGIDENT { // if !reparg && arg == ARGIDENT to replace only once ??
        argv[index] = f
        reparg = true
        continue
      }
      argv[index] = arg
    }
    if !reparg {
      argv = append(argv, f)
    }    
    return exec.Command(command, argv...).Output()
  }
  return newHandler
}

func CreateInternal(command string, fn InternalHandler, args ...string) Handler{
  newHandler := func(f string) ([]byte, error) {
    lcmd := command
    largs := args
    // if this fails, stdin would be empty string
    stdin, _ := SelectStdin()
    return fn(lcmd, stdin, largs...)
  }
  return newHandler  
}
