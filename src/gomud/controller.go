package gomud

const NotifyQueueLength = 100

type Controller interface {
	FailChan() <-chan struct{}
	Notify(Notification)
	Object() DynamicObject
	Start()
	Stop()
}
