package main

import (
	"fmt"
	"os"

	"gomud"
	"time"
)

func innerMain() error {
	// Create two Places, with an interconnecting Edge
	p0 := gomud.NewPlace("The center of the universe; an infinity of light.")
	p1 := gomud.NewPlace("Next to the center of the universe; slightly less light.")
	p1.X = 1
	p1.Y = 1
	_, err := gomud.NewEdge(p0, p1, true, true)
	if err != nil {
		return err
	}

	// Create a Slime with a SimpleAIController
	sState := gomud.SlimeState{Size: 1}
	s := gomud.NewSlime(0, sState)
	s.SetPlace(p0)
	sc, err := gomud.NewSimpleAIController(s, nil)
	if err != nil {
		return err
	}
	sc.Start()

	duration, _ := time.ParseDuration("1s")
	ticker := time.NewTicker(duration)
	i := 1
	for now := range ticker.C {
		fmt.Println("Tick!")
		tick := gomud.TimeTick{When: now}
		s.Controller().Notify(tick)
		if i%3 == 0 {
			b := gomud.Bark{
				Sound: "Woof!",
			}
			s.Controller().Notify(b)
		}
		if i%15 == 0 {
			sc.Stop()
		}
		if i%18 == 0 {
			return nil
		}
		i++
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
