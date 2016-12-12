package model

import (
	"encoding/json"
	"fmt"
	"gomud/supervisor"
	"os"
)

type EventPersister struct {
	eventQueue chan Event
	failChan   chan supervisor.ExitStatus
}

func NewEventPersister() *EventPersister {
	return &EventPersister{
		eventQueue: make(chan Event, EventQueueMaxDepth),
		failChan:   make(chan supervisor.ExitStatus, 0),
	}
}

func (ep *EventPersister) Start() {
	go ep.persistLoop()
}

func (ep *EventPersister) persistLoop() {
	defer func() {
		// FIXME this code is really needed until we have a
		// FIXME supervisor to restart us and report errors/stacktraces
		fmt.Println("OMG EVENTNOTIFIER DIED")
		if r := recover(); r != nil {
			panic(r)
		}
		ep.failChan <- supervisor.FailsafeExit
	}()

	fd, err := os.Create("persist.log")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	for event := range ep.eventQueue {
		fmt.Printf("EventPersister handling: %v\n", event)
		b, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}
		_, err = fmt.Fprint(fd, string(b), "\n")
		if err != nil {
			panic(err)
		}
	}
}
func (ep *EventPersister) EventQueue() chan<- Event {
	return ep.eventQueue
}
