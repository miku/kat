package core

import (  
  "path/filepath"
  "strings"
)

type Handler func(f string) ([]byte, error)
type InternalHandler func(command string, stdin string, args ...string) ([]byte, error)

var moduleTable map[string]Handler 

// default handlers
func init() {
  moduleTable = make(map[string]Handler)
  // -DEBUG-
  // old mapping
  // moduleTable = map[string]Handler{
  //   ".pdf" : PDFHandler,
  //   ".jpg" : ImageHandler,
  //   ".png": ImageHandler,
  //   ".gif": ImageHandler,
  //   ".mrc": ImageHandler,
  //   ".zip": ZipfileHandler,
  //   ".tgz": TarfileHandler,
  //   ".tar.gz": TarfileHandler,
  //   ".mp3": MP3Handler,
  //   ".rar": RarHandler,
  //   ".odt": ODTHandler,
  //   ".docx": ODTHandler,
  //   ".doc":WordHandler,
  //   ".xlsx":XLSXHandler,
  //   ".dmg": DMGHandler,
  //   ".djvu": DjvuHandler,
  //   ".deb": DebianPackageHandler,
  //   ".rpm": RPMHandler,
  // }
  
  for k,v := range defaultHandlers {
    command, args := v[0], v[1:]
    if len(args) == 0 {
      handler := Create(command, "")
      Register(k, handler)      
    } else {
      handler := Create(command, args...)
      Register(k, handler)
    }
  }
}

func Register(extension string, handler Handler) { 
  moduleTable[extension] = handler
}

// Create and register new handler ( similar to calling register with handlers.Create )
func RegisterNewHandler(extension string, args ...string) {
  command, argv := args[0], args[1:]
  Register(extension, Create(command, argv...))
}

// Register internal functions as handlers
func RegisterLocalHandler(command string, fn InternalHandler,args ...string) {
  cmd, argv := args[0], args[1:]
  Register(command, CreateInternal(cmd, fn, argv...))
}

func getHandler(path string, iscmd bool) (Handler, bool) {
  var ext string
  if iscmd {
    ext = path
  } else {
    if ext = filepath.Ext(path); ext != "" {
      p := filepath.Base(path)
      f := strings.SplitN(p, ".", 2)
      ext = "."+f[len(f)-1]
    }    
  }
  handler, status := moduleTable[ext]
  if !status {
    return DefaultHandler, status
  }
  return handler, status
}

func GetHandler(path string) (Handler, bool) { return getHandler(path, false) }
func GetCMDHandler(path string) (Handler, bool) { return getHandler(path, true) }
