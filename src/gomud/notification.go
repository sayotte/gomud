package gomud

import "time"

type Notification interface {
	Status() bool
}

type TimeTick struct {
	When time.Time
}

func (tt TimeTick) Status() bool {
	return true
}

type Bark struct {
	Sound  string
	status bool
}

func (b Bark) Status() bool {
	return b.status
}

type PoisonPill struct{}

func (p PoisonPill) Status() bool {
	return true
}
