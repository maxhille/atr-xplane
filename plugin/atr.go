package atr

import (
	"fmt"
	"sync"
	"time"
)

type Plugin struct {
	api         API
	cmd         chan struct{}
	in          In
	out         Out
	busVolts    DataRef
	fuelPumpOn  DataRef
	fuelPumpOff DataRef
}

type API interface {
	DebugString(string)
	CreateFlightLoop(func() float64) FlightLoopID
	FindDataRef(string) DataRef
	GetDatai(DataRef) int32
	SetDatai(DataRef, int32)
	SetDatavi(DataRef, []int32, int, int)
	GetDatavi(DataRef, int, int) []int32
	GetDatavf(DataRef, int, int) []float32
	GetDataf(DataRef) float32
	RegisterDataAccessor(CustomData) DataRef
	UnregisterDataAccessor(DataRef)
	FindPluginBySignature(string) PluginID
	SendMessageToPlugin(PluginID, PluginMsg, string)
}

type CustomData struct {
	Name       string
	DatatypeID int
	Writable   bool
	Value      []int32
}

const (
	DatatypeUnknown    = 0
	DatatypeInt        = 1
	DatatypeFloat      = 2
	DatatypeDouble     = 4
	DatatypeFloatArray = 8
	DatatypeIntArray   = 16
	DatatypeData       = 32
)

type PluginID int32
type PluginMsg int32

const (
	// https://developer.x-plane.com/code-sample/register-custom-dataref-in-dataref-editor/
	MsgAddDataRef PluginMsg = 0x01000000
	NoPluginID    PluginID  = -1
)

// Opaque types from SDK
type DataRef interface{}
type FlightLoopID interface{}

type input struct {
	busVolts []float32
	fuelPump []int32
}

type output struct {
	fuelPumpOff []int32
}

func New(api API) Plugin {
	return Plugin{
		api: api,
		cmd: make(chan struct{}, 1),
		in:  In{},
		out: Out{},
	}
}

type In struct {
	mu sync.RWMutex
	v  *input
}

type Out struct {
	mu sync.RWMutex
	v  *output
}

func (o *Out) Set(out output) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.v = &out
}

func (i *In) Set(in input) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.v = &in
}

func (o *Out) Get() (*output, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.v, o.v != nil
}

func (i *In) Get() (*input, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.v, i.v != nil
}

func (p *Plugin) Start(name, sig, desc *string) bool {
	*name = "ATR 72-500"
	*desc = "Support plugin for the ATR 72-500"
	*sig = "atr.default"

	p.busVolts = p.api.FindDataRef("sim/cockpit2/electrical/bus_volts")
	p.fuelPumpOn = p.api.FindDataRef("sim/cockpit2/engine/actuators/fuel_pump_on")

	p.fuelPumpOff = p.api.RegisterDataAccessor(CustomData{
		Name:       "atr/fuel_pump_button_off",
		Writable:   true,
		DatatypeID: DatatypeIntArray,
	})

	pluginID := p.api.FindPluginBySignature("xplanesdk.examples.DataRefEditor")
	if pluginID != NoPluginID {
		p.api.SendMessageToPlugin(pluginID, MsgAddDataRef, "atr/fuel_pump_button_off")
	}

	p.startWorker()
	_ = p.api.CreateFlightLoop(p.flightLoop)

	return true
}

func (p *Plugin) Stop() {
	p.stopWorker()
}

func (p *Plugin) flightLoop() float64 {
	busVolts := p.api.GetDatavf(p.busVolts, 0, 2)
	fuelPumpOn := p.api.GetDatavi(p.fuelPumpOn, 0, 2)

	p.in.Set(input{
		busVolts: busVolts,
		fuelPump: fuelPumpOn,
	})

	out, ok := p.out.Get()
	if ok {
		p.api.SetDatavi(p.fuelPumpOff, out.fuelPumpOff, 0, 2)
	}

	return 0.1
}

func (p *Plugin) log(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	p.api.DebugString("ATR 72-500 | " + msg + "\n")
}

func (p *Plugin) startWorker() {
	go func() {
		p.log("starting worker")
		var alive = true
		var ticker = time.NewTicker(100 * time.Millisecond)
		for alive {
			select {
			case <-p.cmd:
				p.log("stopping worker")
				ticker.Stop()
				alive = false
			case <-ticker.C:
				in, ok := p.in.Get()
				if !ok {
					// no value (yet)
					continue
				}
				fuelPumpOff := make([]int32, 2)
				if in.fuelPump[0] == 0 && (in.busVolts[0] > 18.0) {
					fuelPumpOff[0] = 1
				} else {
					fuelPumpOff[0] = 0
				}
				if in.fuelPump[1] == 0 && (in.busVolts[0] > 18.0) {
					fuelPumpOff[1] = 1
				} else {
					fuelPumpOff[1] = 0
				}
				p.out.Set(output{
					fuelPumpOff: fuelPumpOff,
				})
			}
		}
	}()
}

func (p *Plugin) stopWorker() {
	close(p.cmd)
}
