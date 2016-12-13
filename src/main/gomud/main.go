package main

import (
	"fmt"
	"os"

	"gomud"
	"time"
)

func populateExampleWorld(ep *gomud.EventProcessor) error {
	// Create two Places, with an interconnecting Edge
	p0 := gomud.NewPlace(0, "The center of the universe; an infinity of light.")
	ip0 := gomud.InsertPlace{
		NewPlace: p0,
	}
	ep.EventQueue() <- ip0

	p1 := gomud.NewPlace(1, "Next to the center of the universe; slightly less light.")
	p1.X = 1
	p1.Y = 1
	ip1 := gomud.InsertPlace{
		NewPlace: p1,
	}
	ep.EventQueue() <- ip1

	e, err := gomud.NewEdge(p0, p1, true, true)
	if err != nil {
		return err
	}
	ie := gomud.InsertEdge{
		NewEdge: e,
	}
	ep.EventQueue() <- ie

	// Create a Slime
	s := gomud.NewSlime(1000, 1)
	is := gomud.InsertObject{
		NewObject: s,
	}
	ep.EventQueue() <- is

	setPlace := gomud.NewSetPlace(s, p0)
	ep.EventQueue() <- setPlace

	return nil
}

func innerMain() error {
	world := gomud.NewWorld()

	// Create an EventNotifier to broadcast events to concerned
	// objects
	en := gomud.NewEventNotifier(world)
	en.Start()

	// Create an EventPersister to journal all events
	eper := gomud.NewEventPersister()
	eper.Start()

	// Create an EventProcessor to process and persist changes
	// to the world
	ep := gomud.NewSingletonEventProcessor(world, en, eper)
	ep.Start()
	eq := ep.EventQueue()

	// Cheat and load some objects, places, and edges
	err := populateExampleWorld(ep)
	if err != nil {
		return err
	}

	// Start the ticker and let things flow
	duration, _ := time.ParseDuration("1s")
	ticker := time.NewTicker(duration)
	for now := range ticker.C {
		fmt.Println("Tick!")
		eq <- gomud.NewTimeTick(now)
	}
	return nil
}

func main() {
	err := innerMain()
	if err != nil {
		fmt.Printf("Exiting due to error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
