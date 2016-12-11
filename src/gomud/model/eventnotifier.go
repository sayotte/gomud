package model

import (
	"fmt"
	"gomud/supervisor"
)

type EventNotifier struct {
	world      *World
	eventQueue chan Event
	failChan   chan supervisor.ExitStatus
}

func NewEventNotifier(world *World) *EventNotifier {
	return &EventNotifier{
		world:      world,
		eventQueue: make(chan Event, EventQueueMaxDepth),
		failChan:   make(chan supervisor.ExitStatus, 0),
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
func (en *EventNotifier) FailChan() <-chan supervisor.ExitStatus {
	return en.failChan
}

func (en *EventNotifier) notifyLoop() {
	defer func() {
		// FIXME this code is really needed until we have a
		// FIXME supervisor to restart us and report errors/stacktraces
		fmt.Println("OMG EVENTNOTIFIER DIED")
		if r := recover(); r != nil {
			panic(r)
		}
		en.failChan <- supervisor.FailsafeExit
	}()

	for event := range en.eventQueue {
		switch event.(type) {
		case PoisonPill:
			en.failChan <- supervisor.NormalExit
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
