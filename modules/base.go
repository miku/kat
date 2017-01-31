// Extend kat with modules
package modules

import "github.com/miku/kat/core"

type mods struct{
  head *mod
}

type mod struct {
  modname string
  handler core.InternalHandler
  args []string
  next *mod
}

var base *mods = &mods{}

func (self *mods) add(name string, handler core.InternalHandler, args ...string) {
  if self.head == nil {
    self.head = &mod{
      modname: name,
      handler: handler,
      args: args,
      next: nil,
    }
  } else {
    current := self.head
    for ;current.next != nil; current = current.next{}
    current.next = &mod{
      modname: name,
      handler: handler,
      args: args,
      next: nil,
    }
  }  
}

func (self *mods) loadAll() {
  if self.head == nil {
    return
  }
  for current := self.head; current != nil; current = current.next {
    if current.handler == nil {      
      core.RegisterNewHandler(current.modname, current.args...)
    } else {
      core.RegisterLocalHandler(current.modname, current.handler, current.args...)
    }
  }
}

func init() {
  base = &mods{}
}

func loadIntenralModule(name string, handler core.InternalHandler, args ...string) {
  base.add(name, handler, args...)
}

func loadModule(name string, args ...string) {
  base.add(name, nil, args...)
}

// called by initialized ( main ) to load all modules
func Load() {
  base.loadAll()
  base.head.next = nil
  base.head = nil
  base = nil
}
