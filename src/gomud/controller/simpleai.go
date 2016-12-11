package controller

import (
	"fmt"
	"gomud/model"
	"math/rand"
)

func NewSimpleAIController(obj *model.Slime, init interface{}) (Controller, error) {
	saic := &simpleAIController{
		object:     obj,
		failChan:   make(chan struct{}, 0),
		notifyChan: make(chan model.Event, NotifyQueueLength),
		init:       init,
	}
	if err := saic.initialize(); err != nil {
		return nil, err
	}

	return saic, nil
}

type simpleAIController struct {
	object     model.DynamicObject
	failChan   chan struct{}
	notifyChan chan model.Event
	init       interface{}
}

func (saic *simpleAIController) initialize() error {
	saic.object.SetController(saic)
	return nil
}
func (saic *simpleAIController) FailChan() <-chan struct{} {
	return saic.failChan
}
func (saic *simpleAIController) Notify(e model.Event) {
	saic.notifyChan <- e
}
func (saic *simpleAIController) Object() model.DynamicObject {
	return saic.object
}
func (saic *simpleAIController) Start() {
	go saic.controlLoop()
}
func (saic *simpleAIController) Stop() {
	saic.notifyChan <- model.PoisonPill{}
	<-saic.failChan
}
func (saic *simpleAIController) controlLoop() {
	defer func() {
		saic.failChan <- struct{}{}
	}()

	var destination *model.Place
	var edge model.Edge
	var ticksInPlace int
	var ticksToStay int
	ep := model.GetSingletonEventProcessor()
	if ep == nil {
		// If the EventProcessor hasn't been initialized,
		// exit. If we have a supervisor, we'll be restarted.
		fmt.Println("EventProcessor not initialized, SAIC exiting.")
		return
	}

	for notification := range saic.notifyChan {
		fmt.Printf("SAIC: Received notification: %v\n", notification)
		switch notification.(type) {
		case model.TimeTick:
			if saic.object.Place() != nil {
				// We're already in a place if place != nil
				ticksInPlace++
				fmt.Printf("SAIC: ticksInPlace:%d, ticksToStay:%d\n", ticksInPlace, ticksToStay)
				if ticksInPlace > ticksToStay {
					edge, destination = chooseDestination(saic.object.Place())
					setEdge := model.NewSetEdge(saic.object, edge, saic.object.Place())
					ep.EventQueue() <- setEdge
				}
			} else {
				// We're on an edge if place == nil
				setPlace := model.NewSetPlace(saic.object, destination)
				ep.EventQueue() <- setPlace
			}
		case model.PoisonPill:
			return
		case model.SetPlace:
			if notification.ObjectID() == saic.object.ID() {
				fmt.Println("Arriving...")
				ticksInPlace = 0
				ticksToStay = rand.Intn(5)
			} else {
				fmt.Println("Someone else is arriving!")
			}
		case model.SetEdge:
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

func chooseDestination(p *model.Place) (model.Edge, *model.Place) {
	edges := p.Edges()
	for _, dest := range edges[0].OutgoingFromPlaces() {
		if dest != p {
			return edges[0], dest
		}
	}
	panic("Malformed Edge / Place!")
}
