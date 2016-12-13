package gomud

type ModelController interface {
	Notify(Event)
	Stop()
}
