package main

import (
  "github.com/miku/kat/core"
  "github.com/miku/kat/modules"
)

func main() {
  // example: loading registered modules and adding them to handlers
  modules.Load()
  // example: registering new handlers. command, args...
  core.RegisterNewHandler("ip", "dig", "+short", "@resolver1.opendns.com", "myip.opendns.com")
  core.Run()
}




