package controller

import "gomud/model"

const NotifyQueueLength = 100

type Controller interface {
	FailChan() <-chan struct{}
	Notify(model.Event)
	Object() model.DynamicObject
	Start()
	Stop()
}
