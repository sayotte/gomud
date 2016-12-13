package gomud

import (
	"fmt"
	"math/rand"
)

func NewSimpleAIController(obj DynamicObject, init interface{}) (Controller, error) {
	saic := &simpleAIController{
		object:     obj,
		failChan:   make(chan struct{}, 0),
		notifyChan: make(chan Event, NotifyQueueLength),
		init:       init,
	}
	if err := saic.initialize(); err != nil {
		return nil, err
	}

	return saic, nil
}

type simpleAIController struct {
	object     DynamicObject
	failChan   chan struct{}
	notifyChan chan Event
	init       interface{}
}

func (saic *simpleAIController) initialize() error {
	saic.object.SetController(saic)
	return nil
}
func (saic *simpleAIController) FailChan() <-chan struct{} {
	return saic.failChan
}
func (saic *simpleAIController) Notify(e Event) {
	saic.notifyChan <- e
}
func (saic *simpleAIController) Object() DynamicObject {
	return saic.object
}
func (saic *simpleAIController) Start() {
	go saic.controlLoop()
}
func (saic *simpleAIController) Stop() {
	saic.notifyChan <- PoisonPill{}
	<-saic.failChan
}
func (saic *simpleAIController) controlLoop() {
	defer func() {
		saic.failChan <- struct{}{}
	}()

	var destination *Place
	var edge *Edge
	var ticksInPlace int
	var ticksToStay int
	ep := GetSingletonEventProcessor()
	if ep == nil {
		// If the EventProcessor hasn't been initialized,
		// exit. If we have a supervisor, we'll be restarted.
		fmt.Println("EventProcessor not initialized, SAIC exiting.")
		return
	}

	for notification := range saic.notifyChan {
		fmt.Printf("SAIC: Received notification: %v\n", notification)
		switch notification.(type) {
		case TimeTick:
			if saic.object.Place() != nil {
				// We're already in a place if place != nil
				ticksInPlace++
				fmt.Printf("SAIC: ticksInPlace:%d, ticksToStay:%d\n", ticksInPlace, ticksToStay)
				if ticksInPlace > ticksToStay {
					edge, destination = chooseDestination(saic.object.Place())
					setEdge := NewSetEdge(saic.object, edge, saic.object.Place())
					ep.EventQueue() <- setEdge
				}
			} else {
				// We're on an edge if place == nil
				setPlace := NewSetPlace(saic.object, destination)
				ep.EventQueue() <- setPlace
			}
		case PoisonPill:
			return
		case SetPlace:
			if notification.ObjectID() == saic.object.ID() {
				fmt.Println("Arriving...")
				ticksInPlace = 0
				ticksToStay = rand.Intn(5)
			} else {
				fmt.Println("Someone else is arriving!")
			}
		case SetEdge:
			if notification.ObjectID() == saic.object.ID() {
				fmt.Println("Departing BECAUSE I LIKE ANNA")
			} else {
				fmt.Println("Someone else is departing!")
			}
		default:
			fmt.Printf("SAIC not handling this: %v\n", notification)
		}
	}
}

func chooseDestination(p *Place) (*Edge, *Place) {
	edges := p.Edges()
	for _, dest := range edges[0].OutgoingFromPlaces() {
		if dest != p {
			return edges[0], dest
		}
	}
	panic("Malformed Edge / Place!")
}
