package gomud

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

	fdMap := make(map[ObjectID]*os.File)
	defer func(fdMap map[ObjectID]*os.File) {
		for _, fd := range fdMap {
			fd.Close()
		}
	}(fdMap)

	for event := range ep.eventQueue {
		switch event.ObjectID() {
		case NonObjectID:
			continue
		case DoNotRouteID:
			continue
		case BroadcastID:
			continue
		}

		var fd *os.File
		var ok bool
		var err error
		fd, ok = fdMap[event.ObjectID()]
		if !ok {
			fd, err = os.Create(fmt.Sprintf("%v-persist.log", event.ObjectID()))
			if err != nil {
				panic(err)
			}
			fdMap[event.ObjectID()] = fd
		}

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
