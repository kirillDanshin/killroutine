package killroutine

import (
	_ "unsafe"
)

func getg() *g

//go:linkname runtime_suspendG runtime.suspendG
func runtime_suspendG(*g)

//go:linkname goexit0 runtime.goexit0
//go:nosplit
func goexit0(g *g)

//go:linkname casgstatustyped runtime.casgstatus
//go:nosplit
func casgstatustyped(g *g, oldval, newval uint32)

//go:linkname readgstatus runtime.readgstatus
//go:nosplit
func readgstatus(gp *g) uint32

//go:linkname systemstack runtime.systemstack
func systemstack(func())
