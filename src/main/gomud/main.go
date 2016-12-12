package main

import (
	"fmt"
	"os"

	"gomud/controller"
	"gomud/model"
	"time"
)

func genExampleWorld() (*model.World, error) {
	world := model.NewWorld()

	// Create two Places, with an interconnecting Edge
	p0 := model.NewPlace(0, "The center of the universe; an infinity of light.")
	p1 := model.NewPlace(1, "Next to the center of the universe; slightly less light.")
	p1.X = 1
	p1.Y = 1
	e, err := model.NewEdge(p0, p1, true, true)
	if err != nil {
		return nil, err
	}
	world.Places[p0.ID] = p0
	world.Places[p1.ID] = p1
	world.Edges[e.ID()] = e

	// Create a Slime with a SimpleAIController
	s := model.NewSlime(1000, 1)
	world.DynamicObjects[s.ID()] = s
	sc, err := controller.NewSimpleAIController(s, nil)
	if err != nil {
		return nil, err
	}
	sc.Start()

	return world, nil
}

func innerMain() error {
	world, err := genExampleWorld()
	if err != nil {
		return err
	}

	// Create an EventNotifier to broadcast events to concerned
	// objects
	en := model.NewEventNotifier(world)
	en.Start()

	// Create an EventPersister to journal all events
	eper := model.NewEventPersister()
	eper.Start()

	// Create an EventProcessor to process and persist changes
	// to the world
	ep := model.NewSingletonEventProcessor(world, en, eper)
	ep.Start()
	eq := ep.EventQueue()

	setPlace := model.NewSetPlace(world.DynamicObjects[1000], world.Places[0])
	eq <- setPlace

	duration, _ := time.ParseDuration("1s")
	ticker := time.NewTicker(duration)
	for now := range ticker.C {
		fmt.Println("Tick!")
		eq <- model.NewTimeTick(now)
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
