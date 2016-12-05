package gomud

import "fmt"

func NewSimpleAIController(obj *Slime, init interface{}) (Controller, error) {
	saic := &simpleAIController{
		object:     obj,
		failChan:   make(chan struct{}, 0),
		notifyChan: make(chan interface{}, NotifyQueueLength),
		init:       init,
	}
	if err := saic.initialize(); err != nil {
		return nil, err
	}

	return saic, nil
}

type simpleAIController struct {
	object     DynamicObject
	failChan   chan struct{}
	notifyChan chan interface{}
	init       interface{}
}

func (saic *simpleAIController) initialize() error {
	saic.object.SetController(saic)
	return nil
}
func (saic *simpleAIController) FailChan() <-chan struct{} {
	return saic.failChan
}
func (saic *simpleAIController) Notify(n Notification) {
	saic.notifyChan <- n
}
func (saic *simpleAIController) Object() DynamicObject {
	return saic.object
}
func (saic *simpleAIController) Start() {
	go saic.controlLoop()
}
func (saic *simpleAIController) Stop() {
	saic.notifyChan <- PoisonPill{}
	<-saic.failChan
}
func (saic *simpleAIController) controlLoop() {
	defer func() {
		saic.failChan <- struct{}{}
	}()

	for notification := range saic.notifyChan {
		switch notification.(type) {
		case Bark:
			fmt.Println("Heard a bark?")
		case PoisonPill:
			return
		default:
			// do nothing
		}
	}
}
