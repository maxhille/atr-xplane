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
// extern int getDataviCB(void*, int*, int, int);
// extern void setDataviCB(void*, int*, int, int);
import "C"

import (
	"unsafe"

	gopointer "github.com/mattn/go-pointer"

	atr "github.com/maxhille/atr-xplane/plugin"
)

var plugin atr.Plugin
var flcb func() float64

type XPLAPI struct{}

func (api XPLAPI) DebugString(msg string) {
	debugString(msg)
}

func debugString(msg string) {
	cMsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cMsg))
	C.XPLMDebugString(cMsg)
}

//export flightLoopCB
func flightLoopCB(elapsedSinceLastCall, elapsedSinceLastFlightLoop C.float,
	counter C.int, refcon unsafe.Pointer) C.float {
	return C.float(flcb())
}

//export getDataviCB
func getDataviCB(refcon unsafe.Pointer, cVs *C.int, cOffset C.int, cMax C.int) C.int {
	cd := gopointer.Restore(refcon).(*atr.CustomData)

	if cVs == nil {
		return C.int(len(cd.Value))
	}
	ccVs := (*[8]C.int)(unsafe.Pointer(cVs))

	offset := int(cOffset)
	max := int(cMax)
	if max > len(cd.Value)-offset {
		max = len(cd.Value) - offset
	}
	for i := 0; i < max; i++ {
		ccVs[i] = C.int(cd.Value[i+offset])
	}

	// at least datarefeditor seems to need this
	// docs don't though
	return C.int(len(cd.Value))
}

//export setDataviCB
func setDataviCB(refcon unsafe.Pointer, cVs *C.int, cOffset C.int, cCount C.int) {
	cd := gopointer.Restore(refcon).(*atr.CustomData)
	ccVs := (*[8]C.int)(unsafe.Pointer(cVs))

	offset := int(cOffset)
	count := int(cCount)
	cd.Value = make([]int32, count)
	for i := 0; i < count; i++ {
		cd.Value[i+offset] = int32(ccVs[i])
	}
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

func (api XPLAPI) GetDatavf(ref atr.DataRef, offset int, max int) []float32 {
	cRef := ref.(C.XPLMDataRef)
	cVs := [8]C.float{}
	cOffset := C.int(offset)
	cMax := C.int(max)

	C.XPLMGetDatavf(cRef, &cVs[0], cOffset, cMax)

	vs := make([]float32, max)
	for i := 0; i < max; i++ {
		vs[i] = float32(cVs[i])
	}
	return vs
}
func (api XPLAPI) GetDatai(ref atr.DataRef) int32 {
	cRef := ref.(C.XPLMDataRef)
	return int32(C.XPLMGetDatai(cRef))
}

func (api XPLAPI) GetDatavi(ref atr.DataRef, offset int, max int) []int32 {
	cRef := ref.(C.XPLMDataRef)
	cVs := [8]C.int{}
	cOffset := C.int(offset)
	cMax := C.int(max)

	C.XPLMGetDatavi(cRef, &cVs[0], cOffset, cMax)

	vs := make([]int32, max)
	for i := 0; i < max; i++ {
		vs[i] = int32(cVs[i])
	}
	return vs
}

func (api XPLAPI) SetDatai(ref atr.DataRef, v int32) {
	cRef := ref.(C.XPLMDataRef)
	cV := C.int(v)
	C.XPLMSetDatai(cRef, cV)
}

func (api XPLAPI) SetDatavi(ref atr.DataRef, vs []int32, offset int, count int) {
	cRef := ref.(C.XPLMDataRef)
	cOffset := C.int(offset)
	cCount := C.int(count)
	cVs := [8]C.int{}

	if count > len(vs) {
		count = len(vs)
	}
	for i := 0; i < count; i++ {
		cVs[i] = C.int(vs[i])
	}

	C.XPLMSetDatavi(cRef, &cVs[0], cOffset, cCount)
}

func (api XPLAPI) RegisterDataAccessor(cd atr.CustomData) atr.DataRef {
	cName := C.CString(cd.Name)
	defer C.free(unsafe.Pointer(cName))
	cType := C.int(cd.DatatypeID)
	var cWritable C.int
	if cd.Writable {
		cWritable = C.int(1)
	} else {
		cWritable = C.int(0)
	}
	cGetDatavi := C.XPLMGetDatavi_f(unsafe.Pointer(C.getDataviCB))
	cSetDatavi := C.XPLMSetDatavi_f(unsafe.Pointer(C.setDataviCB))
	refcon := gopointer.Save(&cd)

	return C.XPLMRegisterDataAccessor(
		cName,
		cType,
		cWritable,
		nil, nil,
		nil, nil,
		nil, nil,
		cGetDatavi, cSetDatavi,
		nil, nil,
		nil, nil,
		refcon, refcon,
	)
}

func (api XPLAPI) FindPluginBySignature(sig string) atr.PluginID {
	cSig := C.CString(sig)
	defer C.free(unsafe.Pointer(cSig))
	cPluginID := C.XPLMFindPluginBySignature(cSig)
	return atr.PluginID(cPluginID)
}

func (api XPLAPI) SendMessageToPlugin(id atr.PluginID, msg atr.PluginMsg, arg string) {
	cID := C.int(id)
	cMsg := C.int(msg)
	cArg := unsafe.Pointer(C.CString(arg))
	defer C.free(unsafe.Pointer(cArg))
	C.XPLMSendMessageToPlugin(cID, cMsg, cArg)
}

func (api XPLAPI) UnregisterDataAccessor(ref atr.DataRef) {
	// TODO remove gopointer?
	cRef := ref.(C.XPLMDataRef)
	C.XPLMUnregisterDataAccessor(cRef)
}

//export XPluginStart
func XPluginStart(outName, outSig, outDesc *C.char) bool {
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
	plugin.Stop()
}

//export XPluginEnable
func XPluginEnable() int {
	return 1
}

//export XPluginDisable
func XPluginDisable() {
}

//export XPluginReceiveMessage
func XPluginReceiveMessage(inFrom C.XPLMPluginID, inMsg C.int, inParam unsafe.Pointer) {
}

// just for formal reasons
func main() {}
