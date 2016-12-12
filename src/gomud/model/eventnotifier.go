package model

import (
	"fmt"
	"gomud/supervisor"
	"runtime/debug"
)

type EventNotifier struct {
	world      *World
	eventQueue chan Event
	failChan   chan supervisor.PanicWithStack
}

func NewEventNotifier(world *World) *EventNotifier {
	return &EventNotifier{
		world:      world,
		eventQueue: make(chan Event, EventQueueMaxDepth),
		failChan:   make(chan supervisor.PanicWithStack),
	}
}
func (en *EventNotifier) Start() {
	go en.notifyLoop()
}
func (en *EventNotifier) Stop() {
	en.eventQueue <- PoisonPill{}
	<-en.failChan // swallow the supervisor.NormalExit
	<-en.failChan // swallow the supervisor.FailsafeExit
	close(en.failChan)
}
func (en *EventNotifier) FailChan() <-chan supervisor.PanicWithStack {
	return en.failChan
}

func (en *EventNotifier) notifyLoop() {
	defer func() {
		if r := recover(); r != nil {
			pws := supervisor.PanicWithStack{
				PReason: r,
				Stack:   debug.Stack(),
			}
			en.failChan <- pws
			// FIXME here til we're actually supervised
			panic(pws)
		} else {
			close(en.failChan)
		}
	}()

	for event := range en.eventQueue {
		switch event.(type) {
		case PoisonPill:
			return
		case TimeTick:
			en.handleTimeTick(event)
		case SetPlace:
			en.handleSetPlace(event)
		case SetEdge:
			en.handleSetEdge(event)
		default:
			fmt.Printf("EventNotifier not processing this: %v\n", event)
		}
	}
}
func (en *EventNotifier) EventQueue() chan<- Event {
	return en.eventQueue
}
func (en *EventNotifier) handleSetPlace(e Event) {
	sp := e.(SetPlace)
	for _, o := range sp.place.Objects() {
		o.Notify(sp)
	}
}
func (en *EventNotifier) handleTimeTick(e Event) {
	en.world.objLock.RLock()
	for _, o := range en.world.DynamicObjects {
		o.Notify(e)
	}
	en.world.objLock.RUnlock()
}
func (en *EventNotifier) handleSetEdge(e Event) {
	se := e.(SetEdge)
	se.object.Notify(se)
	for _, o := range se.fromPlace.Objects() {
		o.Notify(se)
	}
}
