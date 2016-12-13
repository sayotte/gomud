package gomud

const NotifyQueueLength = 100

type Controller interface {
	FailChan() <-chan struct{}
	Notify(Event)
	Object() DynamicObject
	Start()
	Stop()
}
