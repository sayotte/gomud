package main

import (
	"fmt"
	"os"

	"gomud"
	"time"
)

func innerMain() error {
	sState := gomud.SlimeState{Size: 1}
	s := gomud.NewSlime(0, sState)
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
