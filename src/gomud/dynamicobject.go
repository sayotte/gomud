package gomud

type DynamicObject interface {
	Controller() Controller
	SetController(Controller)
	Location() /*...*/
	Notify(Notification)
	State() DynamicObjectState
}

type DynamicObjectState interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}
