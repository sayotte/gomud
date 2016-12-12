package supervisor

import (
	"reflect"
	"runtime/debug"
)

type PanicWithStack struct {
	PReason interface{}
	Stack   []byte
}

type Supervisable interface {
	Start()
	Stop()
	FailChan() <-chan PanicWithStack
}

type supervisorControlMessage struct {
	msgType string
	payload interface{}
}

type Supervisor struct {
	strategy    string
	children    []Supervisable
	failChan    chan PanicWithStack
	controlChan chan supervisorControlMessage
}

func NewSupervisor(strategy string) *Supervisor {
	s := &Supervisor{
		strategy:    strategy,
		failChan:    make(chan PanicWithStack),
		controlChan: make(chan supervisorControlMessage),
	}
	return s
}

// Supervisor should itself be a Supervisable, so we can make trees
// of them.
func (s *Supervisor) Start() {
	go s.supervisorLoop()
}
func (s *Supervisor) Stop() {
	responseChan := make(chan struct{})
	msg := supervisorControlMessage{
		msgType: "stop",
		payload: responseChan,
	}
	s.controlChan <- msg
	<-responseChan
}
func (s *Supervisor) FailChan() <-chan PanicWithStack {
	return s.failChan
}

func (s *Supervisor) AddChild(c Supervisable) {
	s.controlChan <- supervisorControlMessage{
		msgType: "addchild",
		payload: c,
	}
}

// FIXME this does not stop the child, and that is better
// FIXME done by calling Stop() on the child directly...
// FIXME I'm not sure there's a use-case for this actually,
// FIXME so I'm going to comment it out
/*
func (s *Supervisor) RemoveChild(c Supervisable) {
	s.controlChan <- supervisorControlMessage{
		msgType: "remchild",
		payload: c,
	}
}
*/

func (s *Supervisor) supervisorLoop() {
	defer func() {
		if r := recover(); r != nil {
			pws := PanicWithStack{
				PReason: r,
				Stack:   debug.Stack(),
			}
			s.failChan <- pws
		} else {
			close(s.failChan)
		}
	}()

	for {
		// Add case for the control channel
		selectCases := []reflect.SelectCase{{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(s.controlChan),
		}}
		// Add cases for monitored children
		for _, child := range s.children {
			selectCases = append(selectCases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(child.FailChan()),
			})
		}

		i, recvd, recvdOK := reflect.Select(selectCases)
		if i == 0 {
			// s.controlChan
			msg := recvd.Interface().(supervisorControlMessage)
			switch msg.msgType {
			case "addchild":
				s.children = append(s.children, msg.payload.(Supervisable))
			case "remchild":
				s.removeChildByChild(msg.payload.(Supervisable))
			case "stop":
				// First, stop all our children
				for _, child := range s.children {
					child.Stop()
				}
				// Report successful/synchronous exit
				responseChan := msg.payload.(chan struct{})
				responseChan <- struct{}{}
				return
			}
		} else {
			childIndex := i - 1
			if !recvdOK {
				s.removeChildByIndex(childIndex)
				continue
			}
			// FIXME do something smarter, according to strategy
			s.children[childIndex].Start()
		}

	}
}
func (s *Supervisor) removeChildByChild(c Supervisable) {
	for i, child := range s.children {
		if child == c {
			s.removeChildByIndex(i)
		}
	}
}
func (s *Supervisor) removeChildByIndex(i int) {
	s.children = append(s.children[:i], s.children[i+1:]...)
}
