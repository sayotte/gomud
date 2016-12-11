package model

import (
	"fmt"

	"gomud/supervisor"
)

const EventQueueMaxDepth = 100

var singletonEventProcessor *EventProcessor

type EventProcessor struct {
	world         *World
	eventNotifier *EventNotifier
	eventQueue    chan Event
	failChan      chan supervisor.ExitStatus
}

func NewSingletonEventProcessor(world *World, en *EventNotifier) *EventProcessor {
	singletonEventProcessor = NewEventProcessor(world, en)
	return singletonEventProcessor
}
func GetSingletonEventProcessor() *EventProcessor {
	return singletonEventProcessor
}
func NewEventProcessor(world *World, en *EventNotifier) *EventProcessor {
	ep := &EventProcessor{
		eventNotifier: en,
		world:         world,
	}
	ep.init()
	return ep
}
func (ep *EventProcessor) init() {
	ep.eventQueue = make(chan Event, EventQueueMaxDepth)
	ep.failChan = make(chan supervisor.ExitStatus, 0)
}

// Methods to implement Supervisable interface
func (ep *EventProcessor) Start() {
	go ep.processLoop()
}
func (ep *EventProcessor) Stop() {
	ep.eventQueue <- PoisonPill{}
	<-ep.failChan // swallow the supervisor.NormalExit
	<-ep.failChan // swallow the supervisor.FailsafeExit
	close(ep.failChan)
}
func (ep *EventProcessor) FailChan() <-chan supervisor.ExitStatus {
	return ep.failChan
}

// Functional methods
func (ep *EventProcessor) processLoop() {
	defer func() {
		// FIXME this code is really needed until we have a
		// FIXME supervisor to restart us and report errors/stacktraces
		fmt.Println("OMG EVENTPROCESSOR DIED")
		if r := recover(); r != nil {
			panic(r)
		}
		ep.failChan <- supervisor.FailsafeExit
	}()

	enq := ep.eventNotifier.EventQueue()
	for event := range ep.eventQueue {
		switch event.(type) {
		case SetPlace:
			sp := event.(SetPlace)
			sp.object.setPlace(sp.place)
		case SetEdge:
			se := event.(SetEdge)
			se.object.setEdge(se.edge)
		case PoisonPill:
			ep.failChan <- supervisor.NormalExit
			return
		default:
			// First, journal and replicate
			// Second, decide if it's a change; if so, effect the change
			// Third, pass on to the notification system
			fmt.Printf("EventProcessor not processing this: %v\n", event)
		}
		enq <- event
	}
}
func (ep *EventProcessor) EventQueue() chan<- Event {
	return ep.eventQueue
}
