package xp2ts

import (
	"fmt"
)

type Plugin struct {
	api  API
	conn chan update
	ref  DataRef
}

type API interface {
	DebugString(string)
	CreateFlightLoop(func() float64) FlightLoopID
	FindDataRef(string) DataRef
	GetDatai(DataRef) int32
	GetDataf(DataRef) float32
}

// Opaque types from SDK
type DataRef interface{}
type FlightLoopID interface{}

type update struct {
	v float32
}

func New(api API) Plugin {
	return Plugin{
		api:  api,
		conn: make(chan update, 1),
		ref:  nil,
	}
}

func (p *Plugin) Start(name, sig, desc *string) bool {
	*name = "ATR 72-500"
	*desc = "Support plugin for the ATR 72-500"
	*sig = "atr.default"

	p.ref = p.api.FindDataRef("sim/flightmodel/position/latitude")

	p.startWorker()
	_ = p.api.CreateFlightLoop(p.flightLoop)

	return true
}

func (p *Plugin) Stop() {
	p.stopWorker()
}

func (p *Plugin) flightLoop() float64 {
	v := p.api.GetDataf(p.ref)

	// don't bother the worker if it is doing stuff already
	if len(p.conn) != 0 {
		return 1.0
	}

	p.updateWorker(v)
	return 0.1
}

func (p *Plugin) updateWorker(v float32) {
	p.conn <- update{v}
}

func (p *Plugin) log(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	p.api.DebugString("XP2TS | " + msg + "\n")
}

func (p *Plugin) startWorker() {
	go func() {
		p.log("starting worker")
		for {
			select {
			case _, ok := <-p.conn:
				if !ok {
					p.log("stopping worker")
					return
				}
			}
		}
	}()
}

func (p *Plugin) stopWorker() {
	close(p.conn)
}
