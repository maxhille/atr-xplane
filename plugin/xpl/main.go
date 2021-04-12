package main

// #cgo CFLAGS:-DXPLM200=1 -DXPLM210=1 -DXPLM300=1 -DXPLM301=1 -DXPLM302=1 -DXPLM303=1 -fPIC
// #cgo LDFLAGS: -shared
// #include "XPLM/XPLMInstance.h"
// #include "XPLM/XPLMDisplay.h"
// #include "XPLM/XPLMGraphics.h"
// #include "XPLM/XPLMMenus.h"
// #include "XPLM/XPLMUtilities.h"
// #include "XPLM/XPLMPlugin.h"
// #include "XPLM/XPLMProcessing.h"
// #include "XPLM/XPLMDataAccess.h"
// #include <string.h>
// #include <stdlib.h>
//
// extern float flightLoopCB(float, float, int, void*);
import "C"

import (
	"log"
	"unsafe"

	atr "github.com/maxhille/atr-xplane/plugin"
)

var plugin atr.Plugin
var flcb func() float64

type XPLAPI struct{}

func (api XPLAPI) DebugString(msg string) {
	cMsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cMsg))
	C.XPLMDebugString(cMsg)
}

//export flightLoopCB
func flightLoopCB(elapsedSinceLastCall, elapsedSinceLastFlightLoop C.float,
	counter C.int, refcon unsafe.Pointer) C.float {
	return C.float(flcb())
}

func (api XPLAPI) CreateFlightLoop(cb func() float64) atr.FlightLoopID {
	flcb = cb
	p := C.XPLMCreateFlightLoop_t{
		structSize:   C.sizeof_XPLMCreateFlightLoop_t,
		phase:        C.xplm_FlightLoop_Phase_AfterFlightModel,
		callbackFunc: C.XPLMFlightLoop_f(unsafe.Pointer(C.flightLoopCB)),
		refcon:       nil,
	}
	id := C.XPLMCreateFlightLoop(&p)
	C.XPLMScheduleFlightLoop(id, 1, 0)

	return id
}

func (api XPLAPI) FindDataRef(name string) atr.DataRef {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cRef := C.XPLMFindDataRef(cName)
	return cRef
}

func (api XPLAPI) GetDataf(ref atr.DataRef) float32 {
	cRef := ref.(C.XPLMDataRef)
	return float32(C.XPLMGetDataf(cRef))
}
func (api XPLAPI) GetDatai(ref atr.DataRef) int32 {
	cRef := ref.(C.XPLMDataRef)
	return int32(C.XPLMGetDatai(cRef))
}

//export XPluginStart
func XPluginStart(outName, outSig, outDesc *C.char) bool {
	log.Print("XPluginStart")
	plugin = atr.New(XPLAPI{})
	var name, sig, desc string
	ok := plugin.Start(&name, &sig, &desc)
	goStrCpy(name, outName)
	goStrCpy(sig, outSig)
	goStrCpy(desc, outDesc)

	return ok
}

func goStrCpy(src string, cDst *C.char) {
	cSrc := C.CString(src)
	defer C.free(unsafe.Pointer(cSrc))
	C.strcpy(cDst, cSrc)
}

//export XPluginStop
func XPluginStop() {
	log.Print("XPluginStop")
	plugin.Stop()
}

//export XPluginEnable
func XPluginEnable() int {
	log.Print("XPluginEnable")
	return 1
}

//export XPluginDisable
func XPluginDisable() {
	log.Print("XPluginDisable")
}

//export XPluginReceiveMessage
func XPluginReceiveMessage(inFrom C.XPLMPluginID, inMsg C.int, inParam unsafe.Pointer) {
	log.Print("XPluginReceiveMessage")
}

// just for formal reasons
func main() {}
