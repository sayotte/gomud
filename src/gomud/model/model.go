package model

type ModelController interface {
	Notify(Event)
	Stop()
}
