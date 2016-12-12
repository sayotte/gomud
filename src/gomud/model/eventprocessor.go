package model

import (
	"fmt"

	"gomud/supervisor"
	"runtime/debug"
)

const EventQueueMaxDepth = 100

var singletonEventProcessor *EventProcessor

type EventProcessor struct {
	world          *World
	eventNotifier  *EventNotifier
	eventPersister *EventPersister
	eventQueue     chan Event
	failChan       chan supervisor.PanicWithStack
}

func NewSingletonEventProcessor(world *World, en *EventNotifier, eper *EventPersister) *EventProcessor {
	singletonEventProcessor = NewEventProcessor(world, en, eper)
	return singletonEventProcessor
}
func GetSingletonEventProcessor() *EventProcessor {
	return singletonEventProcessor
}
func NewEventProcessor(world *World, en *EventNotifier, eper *EventPersister) *EventProcessor {
	ep := &EventProcessor{
		eventNotifier:  en,
		eventPersister: eper,
		world:          world,
	}
	ep.init()
	return ep
}
func (ep *EventProcessor) init() {
	ep.eventQueue = make(chan Event, EventQueueMaxDepth)
	ep.failChan = make(chan supervisor.PanicWithStack)
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
func (ep *EventProcessor) FailChan() <-chan supervisor.PanicWithStack {
	return ep.failChan
}

// Functional methods
func (ep *EventProcessor) processLoop() {
	defer func() {
		if r := recover(); r != nil {
			pws := supervisor.PanicWithStack{
				PReason: r,
				Stack:   debug.Stack(),
			}
			ep.failChan <- pws
			// FIXME here til we're actually supervised
			panic(pws)
		} else {
			close(ep.failChan)
		}
	}()

	enq := ep.eventNotifier.EventQueue()
	eperq := ep.eventPersister.EventQueue()
	for event := range ep.eventQueue {
		switch event.(type) {
		case SetPlace:
			sp := event.(SetPlace)
			sp.object.setPlace(sp.place)
		case SetEdge:
			se := event.(SetEdge)
			se.object.setEdge(se.edge)
		case PoisonPill:
			return
		case InsertObject:
			io := event.(InsertObject)
			ep.handleInsertObject(io)
		default:
			// First, journal and replicate
			// Second, decide if it's a change; if so, effect the change
			// Third, pass on to the notification system
			fmt.Printf("EventProcessor not processing this: %v\n", event)
		}
		eperq <- event
		enq <- event
	}
}
func (ep *EventProcessor) EventQueue() chan<- Event {
	return ep.eventQueue
}
func (ep *EventProcessor) handleInsertObject(io InsertObject) {
	return // stub
}
