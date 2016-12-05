package gomud

type DynamicObject interface {
	Controller() Controller
	SetController(Controller)
	Edge() *Edge
	SetEdge(*Edge)
	Place() *Place
	SetPlace(*Place)
	Notify(Notification)
	State() DynamicObjectState
}

type DynamicObjectState interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}
