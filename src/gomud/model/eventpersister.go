package model

import (
	"encoding/json"
	"fmt"
	"gomud/supervisor"
	"os"
	"runtime/debug"
)

type EventPersister struct {
	eventQueue chan Event
	failChan   chan supervisor.PanicWithStack
}

func NewEventPersister() *EventPersister {
	return &EventPersister{
		eventQueue: make(chan Event, EventQueueMaxDepth),
		failChan:   make(chan supervisor.PanicWithStack),
	}
}

func (ep *EventPersister) Start() {
	go ep.persistLoop()
}

func (ep *EventPersister) persistLoop() {
	defer func() {
		if r := recover(); r != nil {
			pws := supervisor.PanicWithStack{
				PReason: r,
				Stack:   debug.Stack(),
			}
			ep.failChan <- pws
			// FIXME here until we're actually supervised
			panic(pws)
		} else {
			close(ep.failChan)
		}
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
